package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
)

func main() {
	opts := immudb.DefaultOptions().WithAddress("localhost").WithPort(3322)
	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(context.TODO(), []byte(`immudb`), []byte(`immudb`), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	defer client.CloseSession(context.TODO())

	idx, err := client.Set(context.TODO(), []byte(`persistedKey`), []byte(`persistedVal`))
	if err != nil {
		log.Fatal(err)
	}

	aOps := &schema.ExecAllRequest{
		Operations: []*schema.Op{
			{
				Operation: &schema.Op_Kv{
					Kv: &schema.KeyValue{
						Key:   []byte(`notPersistedKey`),
						Value: []byte(`notPersistedVal`),
					},
				},
			},
			{
				Operation: &schema.Op_ZAdd{
					ZAdd: &schema.ZAddRequest{
						Set:   []byte(`mySet`),
						Score: 0.4,
						Key:   []byte(`notPersistedKey`)},
				},
			},
			{
				Operation: &schema.Op_ZAdd{
					ZAdd: &schema.ZAddRequest{
						Set:      []byte(`mySet`),
						Score:    0.6,
						Key:      []byte(`persistedKey`),
						AtTx:     idx.Id,
						BoundRef: true,
					},
				},
			},
		},
	}

	idx, err = client.ExecAll(context.TODO(), aOps)
	if err != nil {
		log.Fatal(err)
	}

	list, err := client.ZScan(context.TODO(), &schema.ZScanRequest{
		Set:     []byte(`mySet`),
		SinceTx: idx.Id,
		NoWait:  true,
	})
	if err != nil {
		log.Fatal(err)
	}
	s, _ := json.MarshalIndent(list, "", "\t")
	log.Print(string(s))
}
