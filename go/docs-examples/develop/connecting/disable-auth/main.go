package main

import (
	"context"
	"log"

	immudb "github.com/codenotary/immudb/pkg/client"
)

func main() {
	client, err := immudb.NewImmuClient(
		immudb.DefaultOptions().
			WithAddress("localhost").
			WithPort(3322).
			WithAuth(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.VerifiedSet(context.TODO(), []byte(`immudb`), []byte(`hello world`))
	if err != nil {
		log.Fatal(err)
	}
}
