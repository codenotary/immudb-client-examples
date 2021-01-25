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
	"github.com/codenotary/immudb/pkg/api/schema"
	"log"
	"math"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

func main() {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	idx, _ := client.Set(ctx, []byte(`persistedKey`),[]byte(`persistedVal`))
	_, _ = client.Set(ctx, []byte(`persistedKey`),[]byte(`persistedVal2`))

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

	idx , err = client.ExecAll(ctx, aOps)
	if err != nil {
		log.Fatal(err)
	}
	zscanOpts1 := &schema.ZScanRequest{
		Set:     []byte(`mySet`),
		SinceTx: math.MaxUint64,
		NoWait: true,
	}

	list, err := client.ZScan(ctx, zscanOpts1)
	if err != nil{
		log.Fatal(err)
	}
	s, _ := json.MarshalIndent(list, "", "\t")
	fmt.Print(string(s))
}
