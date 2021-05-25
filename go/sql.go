package main

import (
	"context"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

func main() {

	c, err := client.NewImmuClient(client.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	lr, err := c.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}

	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = c.SQLExec(ctx, `
		BEGIN TRANSACTION
          CREATE TABLE people(id INTEGER, name VARCHAR, salary INTEGER, PRIMARY KEY id);
          CREATE INDEX ON people(name)
		COMMIT
	`, map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.SQLExec(ctx, "UPSERT INTO people(id, name, salary) VALUES (@id, @name, @salary);", map[string]interface{}{"id": 1, "name": "Joe", "salary": 1000})
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.SQLQuery(ctx, "SELECT t.id as d,t.name FROM (people AS t) WHERE id <= 3 AND name = @name", map[string]interface{}{"name": "Joe"}, true)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range res.Rows {
		for _, v := range r.Values {
			log.Printf("%s\n", schema.RenderValue(v.Value))
		}
	}
}
