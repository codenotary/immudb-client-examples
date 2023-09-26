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
	"encoding/json"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
)

// Simple app using official go sdk for immudb

// go mod tidy
// go build
// ./execall

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

	ctx := context.Background()

	idx, err := client.Set(ctx, []byte(`persistedKey`), []byte(`persistedVal`))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Set(ctx, []byte(`persistedKey`), []byte(`persistedVal2`))
	if err != nil {
		log.Fatal(err)
	}

	// Ops payload
	aOps := &schema.ExecAllRequest{
		Operations: []*schema.Op{
			{
				Operation: &schema.Op_Kv{
					Kv: &schema.KeyValue{
						Key:   []byte(`notPersistedKey`),
						Value: []byte(`notPersistedVal`),
					},
				},
			},
			{
				Operation: &schema.Op_ZAdd{
					ZAdd: &schema.ZAddRequest{
						Set:   []byte(`mySet`),
						Score: 0.4,
						Key:   []byte(`notPersistedKey`)},
				},
			},
			{
				Operation: &schema.Op_ZAdd{
					ZAdd: &schema.ZAddRequest{
						Set:      []byte(`mySet`),
						Score:    0.6,
						Key:      []byte(`persistedKey`),
						AtTx:     idx.Id,
						BoundRef: true,
					},
				},
			},
		},
	}
	_, err = client.ExecAll(ctx, aOps)
	if err != nil {
		log.Fatal(err)
	}

	list, err := client.ZScan(ctx, &schema.ZScanRequest{
		Set: []byte(`mySet`),
	})
	if err != nil {
		log.Fatal(err)
	}

	s, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(s))
}
