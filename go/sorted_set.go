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
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
	"log"
	"math"
)

func main() {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	lr , err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil{
		log.Fatal(err)
	}
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	i1, err := client.Set(ctx, []byte(`user1`), []byte(`user1@mail.com`))
	if err != nil{
		log.Fatal(err)
	}
	i2, err := client.Set(ctx, []byte(`user2`), []byte(`user2@mail.com`))
	if err != nil{
		log.Fatal(err)
	}
	i3, err := client.Set(ctx, []byte(`user3`), []byte(`user3@mail.com`))
	if err != nil{
		log.Fatal(err)
	}
	i4, err := client.Set(ctx, []byte(`user3`), []byte(`another-user3@mail.com`))
	if err != nil{
		log.Fatal(err)
	}

	if _ , err = client.ZAddAt(ctx,  []byte(`age`), 25, []byte(`user1`), i1.Id); err != nil {
		log.Fatal(err)
	}
	if _ , err = client.ZAddAt(ctx,  []byte(`age`), 36, []byte(`user2`), i2.Id); err != nil {
		log.Fatal(err)
	}
	if _ , err = client.ZAddAt(ctx,  []byte(`age`), 36, []byte(`user3`), i3.Id); err != nil {
		log.Fatal(err)
	}
	if _ , err = client.ZAddAt(ctx,  []byte(`age`), 54, []byte(`user3`), i4.Id); err != nil {
		log.Fatal(err)
	}

	zscanOpts1 := &schema.ZScanRequest{
		Set:     []byte(`age`),
		SinceTx: math.MaxUint64,
		NoWait: true,
		MinScore: &schema.Score{Score: 36},
	}

	the36YearsOldList, err := client.ZScan(ctx, zscanOpts1)
	if err != nil{
		log.Fatal(err)
	}
	s, _ := json.MarshalIndent(the36YearsOldList, "", "\t")
	fmt.Print(string(s))

	oldestReq := &schema.ZScanRequest{
		Set:           []byte(`age`),
		SeekKey:       []byte{0xFF},
		SeekScore:     math.MaxFloat64,
		SeekAtTx:      math.MaxUint64,
		Limit:         1,
		Desc:          true,
		SinceTx:       math.MaxUint64,
		NoWait:        true,
	}

	oldest, err := client.ZScan(ctx, oldestReq)
	if err != nil{
		log.Fatal(err)
	}
	s, _ = json.MarshalIndent(oldest, "", "\t")
	fmt.Print(string(s))
}
