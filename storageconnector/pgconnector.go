package storageconnector

import (
	"bytes"
	"database/sql"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
	"github.com/zalando-incubator/pazuzu/shared"
)

const (
	createFeaturesTableQuery = `CREATE TABLE IF NOT EXISTS features (
		id serial primary key,
		name TEXT,
		description TEXT,
		author TEXT,
		lastupdate timestamptz,
		dependencies TEXT,
		snippet TEXT,
		test_snippet TEXT
	);`
	getFeatureQuery   = "SELECT * FROM features WHERE name = $1;"
	listFeaturesQuery = "SELECT * FROM features WHERE name ~ $1;"
)

// Reads feature using scanFunc
func readFeature(scanFunc func(dest ...interface{}) error) (shared.Feature, error) {
	var (
		meta         shared.FeatureMeta
		id           int
		dependencies string
		snippet      string
		testSnippet  string
	)

	err := scanFunc(
		&id,
		&meta.Name,
		&meta.Description,
		&meta.Author,
		&meta.UpdatedAt,
		&dependencies,
		&snippet,
		&testSnippet)

	if err != nil {
		return shared.Feature{}, err
	}

	meta.Dependencies = strings.Fields(dependencies)
	buffer := bytes.NewBufferString(testSnippet)

	feature := shared.Feature{
		Meta:        meta,
		Snippet:     snippet,
		TestSnippet: shared.ReadTestSpec(buffer),
	}
	return feature, nil
}

func NewPostgresStorage(connectionString string) (*postgresStorage, error) {
	var pg postgresStorage
	pg.init(connectionString)

	return &pg, nil
}

func createDBConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.Exec(createFeaturesTableQuery)
	return db, nil
}

type postgresStorage struct {
	connectionString string
}

func (store *postgresStorage) init(connectionString string) {
	store.connectionString = connectionString
}

func (store *postgresStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {
	var featureMetas []shared.FeatureMeta

	db, err := createDBConnection(store.connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	featureMetas, err = store.listFeatures(db, name.String())
	if err != nil {
		return featureMetas, err
	}
	return featureMetas, err

}

func (store *postgresStorage) listFeatures(db *sql.DB, name string) ([]shared.FeatureMeta, error) {
	var fms []shared.FeatureMeta

	rows, err := db.Query(listFeaturesQuery, name)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		f, err := readFeature(rows.Scan)
		if err != nil {
			return nil, err
		}
		fms = append(fms, f.Meta)
	}

	return fms, nil
}

func (store *postgresStorage) GetMeta(name string) (shared.FeatureMeta, error) {
	db, err := createDBConnection(store.connectionString)
	if err != nil {
		return shared.FeatureMeta{}, err
	}
	defer db.Close()

	f, err := store.getFeature(db, name)
	if err != nil {
		return shared.FeatureMeta{}, err
	}

	return f.Meta, nil
}

func (store *postgresStorage) GetFeature(name string) (shared.Feature, error) {
	var f shared.Feature

	db, err := createDBConnection(store.connectionString)
	if err != nil {
		return f, err
	}
	defer db.Close()

	return store.getFeature(db, name)
}

func (store *postgresStorage) getFeature(db *sql.DB, name string) (shared.Feature, error) {
	row := db.QueryRow(getFeatureQuery, name)

	f, err := readFeature(row.Scan)
	if err != nil {
		return shared.Feature{}, err
	}

	return f, nil
}

func (store *postgresStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
	var slice []string
	result := map[string]shared.Feature{}
	for _, name := range names {
		err := store.resolve(name, &slice, result)
		if err != nil {
			return []string{}, map[string]shared.Feature{}, err
		}
	}

	return slice, result, nil
}

func (store *postgresStorage) resolve(name string, list *[]string, result map[string]shared.Feature) error {
	if _, ok := result[name]; ok {
		return nil
	}

	feature, err := store.GetFeature(name)
	if err != nil {
		return err
	}
	for _, depName := range feature.Meta.Dependencies {
		err := store.resolve(depName, list, result)
		if err != nil {
			return err
		}
	}

	result[name] = feature
	*list = append(*list, name)
	return nil
}
