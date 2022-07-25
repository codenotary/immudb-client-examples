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

	tx, err := client.Set(context.TODO(), []byte("key1"), []byte("val1"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Set(context.TODO(), []byte("key2"), []byte("val2"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Set(context.TODO(), []byte("key3"), []byte("val3"))
	if err != nil {
		log.Fatal(err)
	}

	txs, err := client.TxScan(context.TODO(), &schema.TxScanRequest{
		InitialTx: tx.Id,
		Limit:     3,
		Desc:      true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Then it's possible to retrieve entries of every transactions:
	for _, tx := range txs.GetTxs() {
		for _, entry := range tx.Entries {
			item, err := client.GetAt(context.TODO(), entry.Key[1:], tx.Header.Id)
			if err != nil {
				item, err = client.GetAt(context.TODO(), entry.Key, tx.Header.Id)
				if err != nil {
					log.Fatal(err)
				}
			}
			log.Printf("retrieved key %s and val %s\n", item.Key, item.Value)
		}
	}
}
