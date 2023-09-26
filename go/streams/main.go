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
	"crypto/rand"
	"io"
	"log"
	"os"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stream"
)

// Simple app using official go sdk for immudb

// go mod tidy
// go build
// ./streams

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

	// first key/value pair with simple values
	key1 := []byte("key1")
	val1 := []byte("val1")

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

	// for the second key/value pair we will put the content of a file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := io.CopyN(tmpfile, rand.Reader, 10*1024); err != nil {
		log.Fatal(err)
	}
	log.Printf("Created temp file with random data: %s", tmpfile.Name())
	tmpfile.Close()

	tmpfile, err = os.Open(tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}

	kv2 := &stream.KeyValue{
		Key: &stream.ValueSize{
			Content: bytes.NewBuffer([]byte(tmpfile.Name())),
			Size:    len(tmpfile.Name()),
		},
		Value: &stream.ValueSize{
			Content: tmpfile,
			Size:    10 * 1024,
		},
	}

	log.Printf("Set values for keys '%s' '%s'", kv1.Key.Content, kv2.Key.Content)
	kvs := []*stream.KeyValue{kv1, kv2}
	_, err = client.StreamSet(ctx, kvs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Reading of key '%s'", tmpfile.Name())

	entry, err := client.StreamGet(ctx, &schema.KeyRequest{Key: []byte(tmpfile.Name())})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("returned key %s with value of len %d", entry.Key, len(entry.Value))

	log.Printf("Chunked reading of key '%s'", tmpfile.Name())

	sc := client.GetServiceClient()
	gs, err := sc.StreamGet(ctx, &schema.KeyRequest{Key: []byte(tmpfile.Name())})
	if err != nil {
		log.Fatal(err)
	}

	kvr := stream.NewKvStreamReceiver(stream.NewMsgReceiver(gs), stream.DefaultChunkSize)

	key, vr, err := kvr.Next()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got value reader for key '%s'", key)

	chunk := make([]byte, 4096)
	for {
		l, err := vr.Read(chunk)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}
		log.Printf("read value chunk: %d byte\n", l)
	}
}
