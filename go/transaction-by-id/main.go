/*
Copyright 2019-2020 vChain, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
)

// Simple app using official go sdk for immudb

// go run main.go

func main() {
	// even though the server address and port are defaults, setting them as a reference
	opts := immudb.DefaultOptions().WithAddress("127.0.0.1").WithPort(3322)

	client := immudb.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := client.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	// ensure connection is closed
	defer client.CloseSession(context.Background())

	// write some entries
	hdr, err := client.SetAll(context.Background(), &schema.SetRequest{
		KVs: []*schema.KeyValue{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
			{
				Key:   []byte("key2"),
				Value: []byte("value2"),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sucessfully commit transaction %d with %d entries\n", hdr.Id, hdr.Nentries)

	// read tx entries as structured values
	tx, err := client.TxByIDWithSpec(context.Background(), &schema.TxRequest{
		Tx: hdr.Id,
		EntriesSpec: &schema.EntriesSpec{
			KvEntriesSpec: &schema.EntryTypeSpec{
				Action: schema.EntryTypeAction_RESOLVE,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range tx.KvEntries {
		fmt.Printf("Got entry: ('%s', '%s')\n", entry.Key, entry.Value)
	}
}
