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

	setTxFirst, err := client.SetAll(context.TODO(),
		&schema.SetRequest{KVs: []*schema.KeyValue{
			{Key: []byte("key1"), Value: []byte("val1")},
			{Key: []byte("key2"), Value: []byte("val2")},
		}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("First txID: %d", setTxFirst.Id)

	// Set keys in another transaction
	setTxSecond, err := client.SetAll(context.TODO(),
		&schema.SetRequest{KVs: []*schema.KeyValue{
			{Key: []byte("key1"), Value: []byte("val11")},
			{Key: []byte("key2"), Value: []byte("val22")},
		}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Second txID: %d", setTxSecond.Id)

	// Without verification
	tx, err := client.TxByID(context.TODO(), setTxFirst.Id)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range tx.Entries {
		item, err := client.GetAt(context.TODO(), entry.Key, setTxFirst.Id)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("retrieved: %+v", item)
	}

	// With verification
	tx, err = client.VerifiedTxByID(context.TODO(), setTxSecond.Id)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range tx.Entries {
		item, err := client.VerifiedGetAt(context.TODO(), entry.Key, setTxSecond.Id)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("retrieved: %+v", item)
	}
}
