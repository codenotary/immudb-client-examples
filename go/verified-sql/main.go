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
	"database/sql"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
)

// Simple app using official go sdk for immudb

// go mod tidy
// go build
// ./app

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

	// Bonus: Here's how you instrument SQLX
	logger := instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
		log.Printf("%s %v", msg, keyvals)
	})

	// Register the new driver as a wrapper of the standard immudb driver
	sql.Register("instrumented-immudb", instrumentedsql.WrapDriver(&stdlib.Driver{}, instrumentedsql.WithLogger(logger)))

	db, err := sqlx.Open("instrumented-immudb", "immudb://immudb:immudb@127.0.0.1:3322/defaultdb?sslmode=disable")
	if err != nil {
		log.Fatal("failed to open DB", err)
	}

	// Create our test table
	_, err = client.SQLExec(context.TODO(), `
		CREATE TABLE IF NOT EXISTS healthchecks_two(
    		id INTEGER AUTO_INCREMENT,
    		name VARCHAR,
   			was_successful boolean NOT NULL,
    		PRIMARY KEY (id)
		);
	`, map[string]interface{}{})
	if err != nil {
		log.Fatal("failed to open DB", err)
	}

	// Insert Data with PGX
	sqlResult := db.MustExec(`INSERT INTO healthchecks_two (name, was_successful) VALUES ($1, $2)`, "Austin", true)

	// PGX gives us the the Last ID, we'll use this to verify our transaction below
	lastId, err := sqlResult.LastInsertId()
	if err != nil {
		log.Fatal(err, " - last insert ID")
	}

	// From here, we're going to construct the arguments for the client.VerifyRow method
	// We'll need
	// 1. The column names - you can get these by querying the data as shown below
	// 	queryResult, err := client.SQLQuery(context.TODO(), `SELECT * FROM healthchecks_two WHERE id = @id`, map[string]interface{}{"id": lastId}, true)
	// 2. The Values we expect to be found in the column with the specified primary key, in the []*schema.SQLValue type
	// 3. The primary key stored as []*schema.SQLValue

	verifyRow := &schema.Row{
		// Here are the column names, as a reminder you can get these out of a queryResult
		Columns: []string{
			"(healthchecks_two.id)",
			"(healthchecks_two.name)",
			"(healthchecks_two.was_successful)",
		},
		// The values we are expecting to find
		Values: []*schema.SQLValue{
			{
				Value: &schema.SQLValue_N{
					N: lastId,
				},
			},
			{
				Value: &schema.SQLValue_S{
					S: "Austin",
				},
			},
			{
				Value: &schema.SQLValue_B{
					B: true,
				},
			},
		},
	}

	// The Primary Key
	PK := []*schema.SQLValue{
		{
			Value: &schema.SQLValue_N{
				N: lastId,
			},
		},
	}

	// Verify the row
	err = client.VerifyRow(context.TODO(), verifyRow, "healthchecks_two", PK)
	if err != nil {
		log.Fatal(err)
	}
}
