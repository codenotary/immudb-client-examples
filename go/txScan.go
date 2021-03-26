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

	_, err = client.Set(ctx, []byte("key1"), []byte("val1"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Set(ctx, []byte("key2"), []byte("val2"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Set(ctx, []byte("key3"), []byte("val3"))
	if err != nil {
		log.Fatal(err)
	}

	txRequest := &schema.TxScanRequest{
		InitialTx: 2,
		Limit:     3,
		Desc:      false,
	}

	txs, err := client.TxScan(ctx, txRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range txs.GetTxs() {
		fmt.Printf("retrieved in ASC tx %d \n", tx.Metadata.Id)
	}
	txRequest = &schema.TxScanRequest{
		InitialTx: 2,
		Limit:     3,
		Desc:      true,
	}

	txs, err = client.TxScan(ctx, txRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range txs.GetTxs() {
		for _, entry := range tx.Entries {
			item, err := client.GetAt(ctx, entry.Key[1:], tx.Metadata.Id)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("retrieved key %s and val %s\n", item.Key, item.Value)
		}
	}

}
