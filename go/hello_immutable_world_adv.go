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

	immuclient "github.com/codenotary/immudb/pkg/client"
)

func main() {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// login with default username and password
	_, err = client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}
	// immudb provides multidatabase capabilities.

	meta, err := client.VerifiedSet(ctx, []byte(`key`), []byte(`val1`))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.VerifiedSet(ctx, []byte(`key`), []byte(`val2`))
	if err != nil {
		log.Fatal(err)
	}

	ventry, err := client.VerifiedGet(ctx, []byte(`key`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully retrieved entry at %v with value %s\n", ventry.Tx, ventry.Value)

	ventry, err = client.VerifiedGetAt(ctx, []byte(`key`), meta.Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully retrieved entry at %v with value %s\n", ventry.Tx, ventry.Value)

	ventry, err = client.VerifiedGetSince(ctx, []byte(`key`), 4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully retrieved entry at %v with value %s\n", ventry.Tx, ventry.Value)
}
