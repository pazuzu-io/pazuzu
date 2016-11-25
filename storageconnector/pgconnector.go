package storageconnector

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"strings"
)

type postgresStorage struct {
	db       *sql.DB
	username string
	dbname   string
}

func (store *postgresStorage) init(username string, dbname string) {
	store.username = username
	store.dbname = dbname
}

func NewPostgresStorage(username string, dbname string) (*postgresStorage, error) {
	var pg postgresStorage
	pg.init(username, dbname)

	return &pg, nil
}

func (store *postgresStorage) connect() error {
	command := fmt.Sprintf("user=%s dbname=%s sslmode=disable", store.username, store.dbname)
	db, err := sql.Open("postgres", command)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	db.Exec("CREATE TABLE IF NOT EXISTS features (index int, name text, description text, author text, lastupdate timestamptz, dependencies text, snippet text);")
	store.db = db
	return nil
}

func (store *postgresStorage) disconnect() {
	store.db.Close()
}

func (store *postgresStorage) scanMeta(SqlQuery string) ([]FeatureMeta, error) {
	var fms []FeatureMeta
	var depText string
	var snippet string
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
		var f FeatureMeta
		err := rows.Scan(&index, &f.Name, &f.Description, &f.Author, &f.UpdatedAt, &depText, &snippet)
		if err != nil {
			return nil, err
		}
		f.Dependencies = strings.Split(depText, " ")
		fms = append(fms, f)
	}
	defer store.disconnect()

	return fms, nil
}

func (store *postgresStorage) SearchMeta(name *regexp.Regexp) ([]FeatureMeta, error) {
	sqlQuery := fmt.Sprintf("select * from features where name ~ '%s';", name)
	fms, err := store.scanMeta(sqlQuery)
	if err != nil {
		return make([]FeatureMeta, 0), err
	}
	return fms, err

}

func (store *postgresStorage) GetMeta(name string) (FeatureMeta, error) {
	sqlQuery := fmt.Sprintf("select * from features where name = '%s';", name)
	fms, err := store.scanMeta(sqlQuery)
	if err != nil {
		return FeatureMeta{}, err
	}

	return fms[0], nil
}

func (store *postgresStorage) GetFeature(name string) (Feature, error) {
	var f Feature
	var index int
	var dep_text string
	sqlQuery := fmt.Sprintf("select * from features where name = '%s';", name)
	store.connect()
	defer store.disconnect()
	err := store.db.QueryRow(sqlQuery).Scan(index, f.Meta.Name, f.Meta.Description, f.Meta.Author, f.Meta.UpdatedAt, dep_text, f.Snippet)
	if err != nil {
		return Feature{}, err
	}

	f.Meta.Dependencies = strings.Split(dep_text, " ")

	return f, nil
}

func (store *postgresStorage) Resolve(names ...string) (map[string]Feature, error) {
	result := map[string]Feature{}
	for _, name := range names {
		err := store.resolve(name, result)
		if err != nil {
			return map[string]Feature{}, err
		}
	}

	return result, nil
}

func (store *postgresStorage) resolve(name string, result map[string]Feature) error {
	if _, ok := result[name]; ok {
		return nil
	}

	feature, err := store.GetFeature(name)
	if err != nil {
		return err
	}

	for _, depName := range feature.Meta.Dependencies {
		err := store.resolve(depName, result)
		if err != nil {
			return err
		}
	}

	result[name] = feature

	return nil
}
