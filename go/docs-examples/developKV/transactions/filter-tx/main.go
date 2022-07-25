package main

import (
	"context"
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

	hdr, err := client.ExecAll(context.TODO(), &schema.ExecAllRequest{
		Operations: []*schema.Op{
			{
				Operation: &schema.Op_Kv{
					Kv: &schema.KeyValue{
						Key:   []byte("key1"),
						Value: []byte("value1"),
					},
				},
			},
			{
				Operation: &schema.Op_Ref{
					Ref: &schema.ReferenceRequest{
						Key:           []byte("ref1"),
						ReferencedKey: []byte("key1"),
					},
				},
			},
			{
				Operation: &schema.Op_ZAdd{
					ZAdd: &schema.ZAddRequest{
						Set:   []byte("set1"),
						Score: 10,
						Key:   []byte("key1"),
					},
				},
			},
		},
	})

	// fetch kv and sorted-set entries as structured values while skipping sql-related entries
	tx, err := client.TxByIDWithSpec(context.TODO(), &schema.TxRequest{
		Tx: hdr.Id,
		EntriesSpec: &schema.EntriesSpec{
			KvEntriesSpec: &schema.EntryTypeSpec{
				Action: schema.EntryTypeAction_RESOLVE,
			},
			ZEntriesSpec: &schema.EntryTypeSpec{
				Action: schema.EntryTypeAction_RESOLVE,
			},
			// explicit exclusion is optional
			SqlEntriesSpec: &schema.EntryTypeSpec{
				// resolution of sql entries is not supported
				Action: schema.EntryTypeAction_EXCLUDE,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range tx.KvEntries {
		log.Printf("retrieved key %s and val %s", entry.Key, entry.Value)
	}

	for _, entry := range tx.ZEntries {
		log.Printf("retrieved set %s key %s and score %v", entry.Set, entry.Key, entry.Score)
	}

	// scan over unresolved entries
	// either EntryTypeAction_ONLY_DIGEST or EntryTypeAction_RAW_VALUE options
	for _, entry := range tx.Entries {
		log.Printf("retrieved key %s and digest %v", entry.Key, entry.HValue)
	}
}
