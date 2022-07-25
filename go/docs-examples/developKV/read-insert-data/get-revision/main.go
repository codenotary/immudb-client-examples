package main

import (
	"context"
	"log"

	immudb "github.com/codenotary/immudb/pkg/client"
)

func main() {
	opts := immudb.DefaultOptions().WithAddress("localhost").WithPort(3322)
	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(context.TODO(), []byte(`immudb`), []byte(`immudb`), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	defer client.CloseSession(context.TODO())

	// Use dedicated API call
	entry, err := client.GetAtRevision(context.TODO(), []byte("key"), -1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Retrieved entry at revision %d: %s", entry.Revision, string(entry.Value))

	// Use additional get option
	entry, err = client.Get(context.TODO(), []byte("key"), immudb.AtRevision(-2))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Retrieved entry at revision %d: %s", entry.Revision, string(entry.Value))
}
