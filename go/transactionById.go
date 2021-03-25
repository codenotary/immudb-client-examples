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
	"github.com/codenotary/immudb/pkg/api/schema"
	"log"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

func main() {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// login with default username and password
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}
	// immudb provides multidatabase capabilities.
	// token is used not only for authentication, but also to route calls to the correct database
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	setRequest := &schema.SetRequest{KVs: []*schema.KeyValue{
		{Key: []byte("key1"), Value: []byte("val1")},
		{Key: []byte("key2"), Value: []byte("val2")},
	}}

	meta, err := client.SetAll(ctx, setRequest)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := client.VerifiedTxByID(ctx, meta.Id)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range tx.Entries {
		item, err := client.VerifiedGetAt(ctx, entry.Key, meta.Id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("retrieved key %s and val %s\n", item.Key, item.Value)
	}
}
