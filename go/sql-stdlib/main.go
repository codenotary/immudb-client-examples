/*
Copyright 2019-2020 vChain, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"fmt"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
)

// Simple app using official go sdk for immudb

// go mod tidy
// go build
// ./sql-stdlib

func main() {
	opts := immudb.DefaultOptions()
	opts.Username = "immudb"
	opts.Password = "immudb"
	opts.Database = "defaultdb"

	db := stdlib.OpenDB(opts)
	defer db.Close()

	_, err := db.ExecContext(context.TODO(), "CREATE TABLE myTable(id INTEGER, name VARCHAR, PRIMARY KEY id)")
	if err != nil {
		panic(err)
	}

	_, err = db.ExecContext(context.TODO(), "INSERT INTO myTable (id, name) VALUES (1, 'immu1')")
	if err != nil {
		panic(err)
	}

	rows, err := db.QueryContext(context.TODO(), "SELECT * FROM myTable")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id uint64
	var name string

	rows.Next()

	err = rows.Scan(&id, &name)
	if err != nil {
		panic(err)
	}

	fmt.Printf("id: %d\n", id)
	fmt.Printf("name: %s\n", name)
}
