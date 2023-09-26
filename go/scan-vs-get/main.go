package main

import (
	"context"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
)

// go mod tidy
// go build
// ./scanning

func main() {
	opts := immudb.DefaultOptions().WithAddress("localhost").WithPort(3322)
	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(context.TODO(), []byte(`immudb`), []byte(`immudb`), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	defer client.CloseSession(context.TODO())

	key := "key1"
	val := "val111"

	hdr, err := client.Set(context.TODO(), []byte(key), []byte(val))
	if err != nil {
		log.Fatal(err)
	}

	entry1, err := client.Get(context.TODO(), []byte(key))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get.  Revision:%d  Key:%s Value:%s \n", entry1.GetRevision(), entry1.GetKey(), entry1.GetValue())

	entries, err := client.Scan(context.TODO(), &schema.ScanRequest{
		Prefix:  []byte(key),
		SinceTx: hdr.Id,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, entry2 := range entries.GetEntries() {
		log.Printf("Scan. Revision:%d  Key:%s Value:%s \n", entry2.GetRevision(), entry2.GetKey(), entry2.GetValue())
	}
}
