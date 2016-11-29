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

const createFeaturesTableQuery = `CREATE TABLE IF NOT EXISTS features (
	id serial primary key,
	name TEXT,
	description TEXT,
	author TEXT,
	lastupdate timestamptz,
	dependencies TEXT,
	snippet TEXT,
	test_snippet TEXT
);`

type PostgresStorage struct {
	db               *sql.DB
	connectionString string
}

func (store *PostgresStorage) init(connectionString string) {
	store.connectionString = connectionString
}

func NewPostgresStorage(connectionString string) (*PostgresStorage, error) {
	var pg PostgresStorage
	pg.init(connectionString)

	return &pg, nil
}

func (store *PostgresStorage) connect() error {
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

func (store *PostgresStorage) disconnect() {
	store.db.Close()
}

func (store *PostgresStorage) scanMeta(SqlQuery string) ([]shared.FeatureMeta, error) {
	var fms []shared.FeatureMeta
	var depText string
	var snippet string
	var testSnippet string
	var index int
	err := store.connect()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(SqlQuery)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f shared.FeatureMeta
		err := rows.Scan(&index, &f.Name, &f.Description, &f.Author, &f.UpdatedAt, &depText, &snippet, &testSnippet)
		if err != nil {
			return nil, err
		}
		f.Dependencies = strings.Fields(depText)
		fms = append(fms, f)
	}
	defer store.disconnect()
	return fms, nil
}

func (store *PostgresStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {
	sqlQuery := fmt.Sprintf("select * from features where name ~ '%s';", name)
	fms, err := store.scanMeta(sqlQuery)
	if err != nil {
		return make([]shared.FeatureMeta, 0), err
	}
	return fms, err

}

func (store *PostgresStorage) GetMeta(name string) (shared.FeatureMeta, error) {
	sqlQuery := fmt.Sprintf("select * from features where name = '%s';", name)
	fms, err := store.scanMeta(sqlQuery)
	if err != nil {
		return shared.FeatureMeta{}, err
	}
	if len(fms) == 0 {
		err = errors.New("Requested feature was not found.")
		return shared.FeatureMeta{}, err
	return fms[0], nil
}

func (store *PostgresStorage) GetFeature(name string) (shared.Feature, error) {
	var f shared.Feature
	var index int
	var dep_text string
	var testSnippet string

	sqlQuery := fmt.Sprintf("select * from features where name = '%s';", name)
	store.connect()
	defer store.disconnect()

	err := store.db.QueryRow(sqlQuery).Scan(&index, &f.Meta.Name, &f.Meta.Description, &f.Meta.Author, &f.Meta.UpdatedAt, &dep_text, &f.Snippet, &testSnippet)
	if err != nil {
		return shared.Feature{}, err
	}

	fmt.Print(f.Snippet)

	buffer := bytes.NewBufferString(testSnippet)

	f.TestSnippet = shared.ReadTestSpec(buffer)
	f.Meta.Dependencies = strings.Fields(dep_text)

	return f, nil
}

func (store *PostgresStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
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

func (store *PostgresStorage) resolve(name string, list *[]string, result map[string]shared.Feature) error {
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
