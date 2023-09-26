package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/codenotary/immudb/pkg/stdlib"
)

func main() {
	// connect with a running immudb server
	db, err := sql.Open("immudb", "immudb://immudb:immudb@127.0.0.1:3322/defaultdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// read some initialization sql script from a file
	bs, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatal(err)
	}

	// execute the sql script
	_, err = db.Exec(string(bs))
	if err != nil {
		log.Fatal(err)
	}

	// lets begin a tx
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// ensure tx is closed (it won't do anything if it's already committed)
	defer tx.Rollback()

	// let's an additional row
	_, err = tx.Exec("INSERT INTO table1(id, title, active) VALUES (3, 'title3', true) ON CONFLICT DO NOTHING;")
	if err != nil {
		log.Fatal(err)
	}

	// commit tx
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// lets query some data
	rows, err := db.Query("SELECT id, title FROM table1 WHERE active = $1", true)
	if err != nil {
		log.Fatal(err)
	}

	var id uint64
	var title string

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &title)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("id: %d, title: %s\n", id, title)
	}
}
