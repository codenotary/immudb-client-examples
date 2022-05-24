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

	// Without verification
	tx, err := client.Set(context.TODO(), []byte(`x`), []byte(`y`))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Set: tx: %d", tx.Id)

	entry, err := client.Get(context.TODO(), []byte(`x`))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get: %v", entry)

	tx, err = client.SetAll(context.TODO(), &schema.SetRequest{
		KVs: []*schema.KeyValue{
			{Key: []byte(`1`), Value: []byte(`test1`)},
			{Key: []byte(`2`), Value: []byte(`test2`)},
			{Key: []byte(`3`), Value: []byte(`test3`)},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("SetAll: tx: %d", tx.Id)

	entries, err := client.GetAll(context.TODO(), [][]byte{[]byte(`1`), []byte(`2`), []byte(`3`)})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GetAll: %+v", entries)

	// With verification
	tx, err = client.VerifiedSet(context.TODO(), []byte(`xx`), []byte(`yy`))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("VerifiedSet: tx: %d", tx.Id)

	entry, err = client.Get(context.TODO(), []byte(`xx`))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("VerifiedGet: %v", entry)
}
