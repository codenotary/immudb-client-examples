package main

import (
	"context"
	"database/sql"
	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"log"
)

func verifySet() {

	// we need both the client as well as a SQL interface
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Login(context.TODO(), []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}

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
CREATE TABLE IF NOT EXISTS healthchecks_two
(
    id INTEGER AUTO_INCREMENT,
    name VARCHAR,
    was_successful boolean NOT NULL,
    PRIMARY KEY (id)
)`, map[string]interface{}{})

	// Insert Data with PGX
	sqlResult := db.MustExec(`
	INSERT INTO healthchecks_two (name, was_successful) VALUES ($1, $2, $3)`, "Austin", true)

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
			"(defaultdb.healthchecks_two.id)",
			"(defaultdb.healthchecks_two.name)",
			"(defaultdb.healthchecks_two.was_successful)",
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
