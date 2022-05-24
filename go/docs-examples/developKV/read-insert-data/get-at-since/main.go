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

	key := []byte(`123123`)
	var txIDs []uint64
	for _, v := range [][]byte{[]byte(`111`), []byte(`222`), []byte(`333`)} {
		txID, err := client.Set(context.TODO(), key, v)
		if err != nil {
			log.Fatal(err)
		}
		txIDs = append(txIDs, txID.Id)
	}

	otherTxID, err := client.Set(context.TODO(), []byte(`other`), []byte(`other`))
	if err != nil {
		log.Fatal(err)
	}

	// Without verification
	entry, err := client.GetSince(context.TODO(), key, txIDs[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GetSince first: %+v", entry)

	// With verification
	entry, err = client.VerifiedGetSince(context.TODO(), key, txIDs[0]+1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("VerifiedGetSince second: %+v", entry)

	// GetAt txID after inserting other data
	_, err = client.GetAt(context.TODO(), key, otherTxID.Id)
	if err == nil {
		log.Fatalf("This should not happen, %+v", entry)
	}

	// Without verification
	entry, err = client.GetAt(context.TODO(), key, txIDs[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GetAt second: %+v", entry)

	// With verification
	entry, err = client.VerifiedGetAt(context.TODO(), key, txIDs[2])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("VerifiedGetAt third: %+v", entry)

	// VerifiedGetAt txID after inserting other data
	entry, err = client.VerifiedGetAt(context.TODO(), key, otherTxID.Id)
	if err == nil {
		log.Fatalf("This should not happen, %+v", entry)
	}
}
