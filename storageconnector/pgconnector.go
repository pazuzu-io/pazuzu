package storageconnector

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
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
	getFeatureQuery    = "SELECT * FROM features WHERE name = '%s';"
	searchFeatureQuery = "SELECT * FROM features WHERE name ~ '%s';"
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

type postgresStorage struct {
	db               *sql.DB
	connectionString string
}

func (store *postgresStorage) init(connectionString string) {
	store.connectionString = connectionString
}

func (store *postgresStorage) connect() error {
	db, err := sql.Open("postgres", store.connectionString)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	db.Exec(createFeaturesTableQuery)
	store.db = db
	return nil
}

func (store *postgresStorage) disconnect() {
	store.db.Close()
}

func (store *postgresStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {
	sqlQuery := fmt.Sprintf(searchFeatureQuery, name)
	fms, err := store.scanMeta(sqlQuery)
	if err != nil {
		return make([]shared.FeatureMeta, 0), err
	}
	return fms, err

}

func (store *postgresStorage) scanMeta(SqlQuery string) ([]shared.FeatureMeta, error) {
	var fms []shared.FeatureMeta

	err := store.connect()
	if err != nil {
		return nil, err
	}
	defer store.disconnect()

	rows, err := store.db.Query(SqlQuery)
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
	sqlQuery := fmt.Sprintf(getFeatureQuery, name)
	fms, err := store.scanMeta(sqlQuery)
	if err != nil {
		return shared.FeatureMeta{}, err
	}
	if len(fms) == 0 {
		err = errors.New("Requested feature was not found.")
		return shared.FeatureMeta{}, err
	}
	return fms[0], nil
}

func (store *postgresStorage) GetFeature(name string) (shared.Feature, error) {
	var f shared.Feature

	sqlQuery := fmt.Sprintf(getFeatureQuery, name)
	store.connect()
	defer store.disconnect()

	row := store.db.QueryRow(sqlQuery)
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
