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

	_, err = client.Set(context.Background(), []byte("user1"), []byte("user1@mail.com"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Set(context.Background(), []byte("user2"), []byte("user2@mail.com"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Set(context.Background(), []byte("user3"), []byte("user3@mail.com"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Set(context.Background(), []byte("user3"), []byte("another-user3@mail.com"))
	if err != nil {
		log.Fatal(err)
	}

	if _, err = client.ZAdd(context.Background(), []byte("age"), 25, []byte("user1")); err != nil {
		log.Fatal(err)
	}
	if _, err = client.ZAdd(context.Background(), []byte("age"), 50, []byte("user2")); err != nil {
		log.Fatal(err)
	}
	if _, err = client.ZAdd(context.Background(), []byte("age"), 36, []byte("user3")); err != nil {
		log.Fatal(err)
	}

	resp, err := client.ZScan(context.Background(), &schema.ZScanRequest{
		Set:      []byte("age"),
		MinScore: &schema.Score{Score: 30},
	})
	if err != nil {
		log.Fatal(err)
	}

	s, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Printf("older than %d: %s\n", 30, s)

	fmt.Println()

	oldest, err := client.ZScan(context.Background(), &schema.ZScanRequest{
		Set:   []byte("age"),
		Desc:  true,
		Limit: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	s, _ = json.MarshalIndent(oldest, "", "\t")
	fmt.Printf("oldest: %s\n", s)
}
