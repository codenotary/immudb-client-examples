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
	"github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	c, err := client.NewImmuClient(client.DefaultOptions().WithServerSigningPubKey("../example-public.key"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	lr, err := c.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}

	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	if _, err := c.Set(ctx, []byte(`immudb`), []byte(`hello world`)); err != nil {
		log.Fatal(err)
	}

	var state *schema.ImmutableState
	if state, err = c.CurrentState(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Print(state)
}
