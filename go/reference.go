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
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	_, err = client.Set(ctx, []byte(`firstKey`), []byte(`firstValue`))
	if err != nil {
		log.Fatal(err)
	}
	reference, err := client.SetReference(ctx, []byte(`myTag`), []byte(`firstKey`))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", reference)
	firstItem, err := client.Get(ctx, []byte(`myTag`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", firstItem)

	_, err = client.Set(ctx, []byte(`secondKey`), []byte(`secondValue`))
	if err != nil {
		log.Fatal(err)
	}
	reference, err = client.VerifiedSetReference(ctx, []byte(`mySecondTag`), []byte(`secondKey`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", reference)

	secondItem, err := client.Get(ctx, []byte(`mySecondTag`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", secondItem)

	meta, err := client.Set(ctx, []byte(`secondKey`), []byte(`secondValue`))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Set(ctx, []byte(`secondKey`), []byte(`thirdValue`))
	if err != nil {
		log.Fatal(err)
	}
	reference, err = client.VerifiedSetReferenceAt(ctx, []byte(`myThirdTag`), []byte(`secondKey`), meta.Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", reference)

	_, err = client.Set(ctx, []byte(`secondKey`), []byte(`secondValue`))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Set(ctx, []byte(`secondKey`), []byte(`thirdValue`))
	if err != nil {
		log.Fatal(err)
	}
	reference, err = client.VerifiedSetReference(ctx, []byte(`myThirdTag`), []byte(`secondKey`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", reference)

	thirdItem, err := client.Get(ctx, []byte(`myThirdTag`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", thirdItem)
}
