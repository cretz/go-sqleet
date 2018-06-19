package main

import (
	"database/sql"
	"flag"
	"log"
	"net/url"
	"strings"

	_ "github.com/cretz/go-sqleet/sqlite3"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbPath := "simple.db"
	flag.StringVar(&dbPath, "db", dbPath, "Path to db file")
	key := ""
	flag.StringVar(&key, "key", key, "Key used to open/set")
	rekey := ""
	flag.StringVar(&rekey, "rekey", rekey, "Key used to reset on write")
	flag.Parse()

	// Build the data source name with parameters
	dsn := dbPath
	if key != "" {
		log.Printf("Opening/creating with a key")
		dsn += "?_key=" + url.QueryEscape(key)
	} else {
		log.Printf("Not setting any encryption key")
	}
	if rekey != "" {
		log.Printf("Changing with a rekey")
		if strings.Contains(dsn, "?") {
			dsn += "&"
		} else {
			dsn += "?"
		}
		dsn += "_rekey=" + url.QueryEscape(rekey)
	}

	// Open the DB
	log.Printf("Opening/creating DB %v", dsn)
	db, err := sql.Open("sqleet", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Printf("Creating table 'foo' if not present")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)")
	if err != nil {
		return err
	}

	log.Printf("Adding 'bar', 'baz', and 'qux' name values")
	_, err = db.Exec("INSERT INTO foo (name) VALUES (?), (?), (?)", "bar", "baz", "qux")
	if err != nil {
		return err
	}

	log.Printf("Fetching all name values from 'foo':")
	rows, err := db.Query("SELECT name FROM foo")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return err
		}
		log.Printf("Name: %v", name)
	}
	return nil
}
