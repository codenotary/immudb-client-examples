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

	_, err = client.Set(context.TODO(), []byte(`key1`), []byte(`val1`))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Set(context.TODO(), []byte(`key2`), []byte(`val2`))
	if err != nil {
		log.Fatal(err)
	}

	itList, err := client.GetAll(context.TODO(), [][]byte{
		[]byte("key1"),
		[]byte("key2"),
		[]byte("key3"), // does not exist, no value returned
	})

	log.Printf("Set: tx: %+v", itList)
}
