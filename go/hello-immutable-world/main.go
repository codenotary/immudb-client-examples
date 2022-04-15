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

	immudb "github.com/codenotary/immudb/pkg/client"
)

// Simple app using official go sdk for immudb

// go run main.go

func main() {
	// even though the server address and port are defaults, setting them as a reference
	opts := immudb.DefaultOptions().WithAddress("127.0.0.1").WithPort(3322)

	client := immudb.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := client.OpenSession(context.Background(), []byte(`immudb`), []byte(`immudb`), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	// ensure connection is closed
	defer client.CloseSession(context.Background())

	// write an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	hdr, err := client.VerifiedSet(context.Background(), []byte(`hello`), []byte(`immutable world`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sucessfully set a verified entry: ('%s', '%s') @ tx %d\n", []byte(`hello`), []byte(`immutable world`), hdr.Id)

	// read an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	entry, err := client.VerifiedGet(context.Background(), []byte(`hello`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sucessfully got veriified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)
}
