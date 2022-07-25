package main

import (
	"context"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
)

func main() {
	opts := immudb.DefaultOptions().
		WithAddress("localhost").
		WithPort(3322)

	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(context.TODO(), []byte(`immudb`), []byte(`immudb`), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	// To perform SQL statements, use the SQLExec function, which takes a SQLExecRequest with a SQL operation:
	_, err = client.SQLExec(context.TODO(), `
	    BEGIN TRANSACTION;
	            CREATE TABLE people(id INTEGER, name VARCHAR[256], salary INTEGER, PRIMARY KEY id);
	            CREATE INDEX ON people(name);
	    COMMIT;
	    `,
		map[string]interface{}{},
	)
	if err != nil {
		log.Fatal(err)
	}

	// This is also how you perform inserts:
	_, err = client.SQLExec(context.TODO(),
		"UPSERT INTO people(id, name, salary) VALUES (@id, @name, @salary);",
		map[string]interface{}{"id": 1, "name": "Joe", "salary": 1000},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Once you have data in the database, you can use the SQLQuery method of the client to query.
	// Both SQLQuery and SQLExec allows named parameters.
	// Just encode them as @param and pass map[string]{}interface as values:
	res, err := client.SQLQuery(context.TODO(),
		"SELECT t.id AS d, t.name FROM people AS t WHERE id <= 3 AND name = @name",
		map[string]interface{}{"name": "Joe"},
		true,
	)
	if err != nil {
		fmt.Printf("To tu ?")
		log.Fatal(err)
	}

	// res is of the type *schema.SQLQueryResult. In order to iterate over the results,
	// you iterate over res.Rows. On each iteration, the row r will have a member Values,
	// which you can iterate to get each column.
	for _, r := range res.Rows {
		for _, v := range r.Values {
			log.Printf("%s\n", schema.RenderValue(v.Value))
		}
	}
}
