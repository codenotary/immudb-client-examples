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
	"log"

	immudb "github.com/codenotary/immudb/pkg/client"
)

// Simple app using official go sdk for immudb

// go run main.go

func main() {
	// even though the server address and port are defaults, setting them as a reference
	opts := immudb.DefaultOptions().WithAddress("127.0.0.1").WithPort(3322)

	client := immudb.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := client.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	// ensure connection is closed
	defer client.CloseSession(context.Background())

	_, err = client.SQLExec(context.Background(), `
		BEGIN TRANSACTION;
          CREATE TABLE IF NOT EXISTS people(id INTEGER, name VARCHAR[50], salary INTEGER, PRIMARY KEY id);
          CREATE INDEX IF NOT EXISTS ON people(name);
		COMMIT;
	`, map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sucessfully created table and index\n")

	_, err = client.SQLExec(context.Background(), "UPSERT INTO people(id, name, salary) VALUES (@id, @name, @salary);", map[string]interface{}{"id": 1, "name": "Joe", "salary": 1000})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sucessfully row insertion\n")

	res, err := client.SQLQuery(context.Background(), "SELECT t.id as d,t.name FROM people AS t WHERE id <= 3 AND name = @name", map[string]interface{}{"name": "Joe"}, true)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range res.Rows {
		log.Printf("Got row: %v\n", row)
	}
}
