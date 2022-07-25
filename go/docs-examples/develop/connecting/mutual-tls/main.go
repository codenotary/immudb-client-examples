package main

import (
	"context"
	"log"

	immudb "github.com/codenotary/immudb/pkg/client"
)

func main() {
	// Folder cotaining MTLS certificates
	pathToMTLSFolder := "./mtls"

	opts := immudb.DefaultOptions().
		WithAddress("localhost").
		WithPort(3322).
		WithMTLs(true).
		WithMTLsOptions(
			immudb.MTLsOptions{}.
				WithCertificate(pathToMTLSFolder + "/4_client/certs/localhost.cert.pem").
				WithPkey(pathToMTLSFolder + "/4_client/private/localhost.key.pem").
				WithClientCAs(pathToMTLSFolder + "/2_intermediate/certs/ca-chain.cert.pem").
				WithServername("localhost"),
		)

	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(context.TODO(), []byte(`immudb`), []byte(`immudb`), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	defer client.CloseSession(context.TODO())

	// do amazing stuff
}
