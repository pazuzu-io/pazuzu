package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
)

type postgreStorage struct {
	db *sql.DB
}

func ConnectToPg (username string, dbname string) {
	fmt.Println("Opening database connection.")
	command := fmt.Sprintf("user=%s dbname=%s sslmode=disable",username,dbname)
	db, err := sql.Open("postgres", command)

	checkErr(err)
	fmt.Println("Successfully opened database connection.")
	db.Exec("CREATE TABLE IF NOT EXISTS features (index int, name text, description text, author text, lastupdate timestamptz, dependencies text);")
}

func (store* postgreStorage) disconnectFromPg() {
	store.db.Close()
}

func (store* postgreStorage) SearchMeta(name *regexp.Regexp) ([]FeatureMeta, error) {
	return nil, nil
}

func (store* postgreStorage) GetMeta(name string) (FeatureMeta, error) {
	return nil, nil
}

func (store* postgreStorage) GetFeature(name string) (Feature, error) {
	return nil,nil
}

func (store* postgreStorage) Resolve(names ...string) (map[string]Feature, error) {
	return nil,nil
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
