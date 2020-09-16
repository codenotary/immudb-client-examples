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
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// login with default username and password
	lr , err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))

	// immudb provides multidatabase capabilities.
	// token is used not only for authentication, but also to route calls to the correct database
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	if _, err := client.Set(ctx, []byte(`immudb`), []byte(`hello world`)); err != nil {
		log.Fatal(err)
	}

	if item, err := client.Get(ctx, []byte(`immudb`)); err != nil {
		log.Fatal(err)
	}else{
		// immudb sdk provides structured data. https://github.com/codenotary/immudb#structured-value
		fmt.Printf("%s\n", item.Value.Payload)
	}
}
