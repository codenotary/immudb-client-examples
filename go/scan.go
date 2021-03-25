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

	_, _ = client.Set(ctx, []byte(`aaa`), []byte(`item1`))
	_, _ = client.Set(ctx, []byte(`bbb`), []byte(`item2`))
	_, _ = client.Set(ctx, []byte(`abc`), []byte(`item3`))

	scanReq := &schema.ScanRequest{
		Prefix: []byte(`a`),
	}

	list, err := client.Scan(ctx, scanReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", list)
	scanReq1 := &schema.ScanRequest{
		SeekKey: []byte{0xFF},
		Prefix:  []byte(`a`),
		Desc:    true,
	}

	list, err = client.Scan(ctx, scanReq1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", list)
	scanReq2 := &schema.ScanRequest{
		SeekKey: []byte{0xFF},
		Desc:    true,
		SinceTx: math.MaxUint64,
		NoWait:  true,
	}

	list, err = client.Scan(ctx, scanReq2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", list)
}
