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
	"bytes"
	"context"
	"fmt"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/stream"
	"io"
	"log"
	"os"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

func main() {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions().WithStreamChunkSize(4096))
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

	myFileName := "/tmp/libunix-java13791202755704945675so"
	key1 := []byte("key1")
	val1 := []byte("val1")
	f, err := os.Open(myFileName)
	if err != nil {
		log.Fatal(err)
	}
	stats, err := os.Stat(myFileName)
	if err != nil {
		log.Fatal(err)
	}

	kv1 := &stream.KeyValue{
		Key: &stream.ValueSize{
			Content: bytes.NewBuffer(key1),
			Size:    len(key1),
		},
		Value: &stream.ValueSize{
			Content: bytes.NewBuffer(val1),
			Size:    len(val1),
		},
	}
	kv2 := &stream.KeyValue{
		Key: &stream.ValueSize{
			Content: bytes.NewBuffer([]byte(myFileName)),
			Size:    len(myFileName),
		},
		Value: &stream.ValueSize{
			Content: f,
			Size:    int(stats.Size()),
		},
	}

	kvs := []*stream.KeyValue{kv1, kv2}
	_, err = client.StreamSet(ctx, kvs)
	if err != nil {
		log.Fatal(err)
	}

	entry, err := client.StreamGet(ctx, &schema.KeyRequest{ Key: []byte(myFileName)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("returned key %s", entry.Key)

	sc := client.GetServiceClient()
	gs, err := sc.StreamGet(ctx, &schema.KeyRequest{ Key: []byte(myFileName)})

	kvr := stream.NewKvStreamReceiver(stream.NewMsgReceiver(gs), stream.DefaultChunkSize)

	key, vr, err := kvr.Next()
	fmt.Printf("read %s key", key)
	if err != nil {
		log.Fatal(err)
	}

	chunk := make([]byte, 4096)
	for {
		l, err := vr.Read(chunk)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Printf("read %d byte\n", l)
	}


}
