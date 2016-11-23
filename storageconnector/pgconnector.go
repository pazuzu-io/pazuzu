package storageconnector

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"strings"
)

type postgreStorage struct {
	db *sql.DB
}

func ConnectToPg(username string, dbname string) {
	fmt.Println("Opening database connection.")
	command := fmt.Sprintf("user=%s dbname=%s sslmode=disable", username, dbname)
	db, err := sql.Open("postgres", command)

	checkErr(err)
	fmt.Println("Successfully opened database connection.")
	db.Exec("CREATE TABLE IF NOT EXISTS features (index int, name text, description text, author text, lastupdate timestamptz, dependencies text);")
}

func (store *postgreStorage) disconnectFromPg() {
	store.db.Close()
}

func (store *postgreStorage) SearchMeta(name *regexp.Regexp) ([]FeatureMeta, error) {
	// TODO: implement (github issue #91)
	return nil, nil
}

func (store *postgreStorage) GetMeta(name string) (FeatureMeta, error) {
	var f FeatureMeta
	var sql_query string
	var index int
	var dep_text string
	sql_query = fmt.Sprintf("select * from features where name == %s", name)
	err := store.db.QueryRow(sql_query).Scan(index, f.Name, f.Description, f.Author, f.UpdatedAt, dep_text)

	checkErr(err)

	s := strings.Split(dep_text, " ")
	f.Dependencies = make([]string, len(s))

	for index, value := range s {
		f.Dependencies[index] = value
	}
	return f, nil
}

func (store *postgreStorage) GetFeature(name string) (Feature, error) {
	// TODO: implement (github issue #91)
	var f Feature
	return f, nil
}

func (store *postgreStorage) Resolve(names ...string) (map[string]Feature, error) {
	// TODO: implement (github issue #91)
	return nil, nil
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
