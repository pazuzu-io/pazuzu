package storageconnector

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"strings"
)

type postgreStorage struct {
	db       *sql.DB
	username string
	dbname   string
}

func (store *postgreStorage) init(username string, dbname string) {
	store.username = username
	store.dbname = dbname
}

func (store *postgreStorage) connect() {
	command := fmt.Sprintf("user=%s dbname=%s sslmode=disable", store.username, store.dbname)
	db, err := sql.Open("postgres", command)

	checkErr(err)
	db.Exec("CREATE TABLE IF NOT EXISTS features (index int, name text, description text, author text, lastupdate timestamptz, dependencies text, snippet text);")
}

func (store *postgreStorage) disconnect() {
	store.db.Close()
}

func (store *postgreStorage) scanMeta(SqlQuery string) ([]FeatureMeta, error) {
	var fms []FeatureMeta
	var depText string
	var snippet string
	var index int
	store.connect()
	defer store.disconnect()
	rows, err := store.db.Query(SqlQuery)
	checkErr(err)
	for rows.Next() {
		var f FeatureMeta
		err := rows.Scan(index, f.Name, f.Description, f.Author, f.UpdatedAt, depText, snippet)
		checkErr(err)

		f.Dependencies = strings.Split(depText, " ")
		fms = append(fms, f)
	}

	return fms, nil
}

func (store *postgreStorage) SearchMeta(name *regexp.Regexp) ([]FeatureMeta, error) {
	sqlQuery := fmt.Sprintf("select * from features where name ~ %s", name)
	fms, err := store.scanMeta(sqlQuery)
	checkErr(err)
	return fms, err

}

func (store *postgreStorage) GetMeta(name string) (FeatureMeta, error) {
	sqlQuery := fmt.Sprintf("select * from features where name == %s", name)
	fms, err := store.scanMeta(sqlQuery)

	checkErr(err)

	return fms[0], nil
}

func (store *postgreStorage) GetFeature(name string) (Feature, error) {
	var f Feature
	var index int
	var dep_text string
	sqlQuery := fmt.Sprintf("select * from features where name == %s", name)
	err := store.db.QueryRow(sqlQuery).Scan(index, f.Meta.Name, f.Meta.Description, f.Meta.Author, f.Meta.UpdatedAt, dep_text, f.Snippet)

	checkErr(err)

	f.Meta.Dependencies = strings.Split(dep_text, " ")

	return f, nil
}

func (store *postgreStorage) Resolve(names ...string) (map[string]Feature, error) {
	result := map[string]Feature{}
	for _, name := range names {
		err := store.resolve(name, result)
		if err != nil {
			return map[string]Feature{}, err
		}
	}

	return result, nil
}

func (store *postgreStorage) resolve(name string, result map[string]Feature) error {
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

// main function for testing purposes only
func main() {
	fmt.Println("Opening database connection.")
	db, err := sql.Open("postgres", "user=omaurer dbname=pazuzu sslmode=disable")

	checkErr(err)
	fmt.Println("Successfully opened database connection.")

	defer db.Close()

	var index int
	err = db.QueryRow(`INSERT INTO features values (1, 'Olaf', 'super', 'nochmehr', '2003-04-12 04:05:06+02', 'java python');`).Scan(&index)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error error")
		panic(err)
	}
}
