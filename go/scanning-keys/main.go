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
	for _, keyPrefix := range []string{"a", "b", "c"} {
		for i := 0; i < 3; i++ {
			key := []byte(fmt.Sprintf("%s_%d", keyPrefix, i))
			value := []byte(fmt.Sprintf("val_%s_%d", keyPrefix, i))

			_, err = client.Set(context.Background(), key, value)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// scan over all keys with prefix "b" i.e. "b_0", "b_1" and "b_2"
	resp, err := client.Scan(context.Background(), &schema.ScanRequest{
		Prefix: []byte("b"),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range resp.Entries {
		fmt.Printf("Got entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)
	}
}
