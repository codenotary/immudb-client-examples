package main

import (
	"context"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
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

	tx, err := client.SetAll(context.TODO(), &schema.SetRequest{
		KVs: []*schema.KeyValue{
			{Key: []byte(`1`), Value: []byte(`key1`)},
			{Key: []byte(`2`), Value: []byte(`key2`)},
			{Key: []byte(`3`), Value: []byte(`key3`)},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("SetAll: tx: %d", tx.Id)
}
