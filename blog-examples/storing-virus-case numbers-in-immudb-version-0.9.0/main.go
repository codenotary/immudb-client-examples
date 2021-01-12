package main
import (
	"context"
	"fmt"
	"log"
)

import (
	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)




func main() {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// login with default username and password
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}

	// immudb provides multidatabase capabilities.
	// token is used not only for authentication, but also to route calls to the correct database
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

		// creating new database
		err = client.CreateDatabase(ctx, &schema.Database{
			Databasename: "antigentestdb",
		})
		if err != nil {
			log.Fatal(err)
		}

	// switch to database
	resp, err := client.UseDatabase(ctx, &schema.Database{
		Databasename: "antigentestdb",
	})

	md = metadata.Pairs("authorization", resp.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// create example data
	var antigenData = struct {
		keys    [][]byte
		values  [][]byte
		refKeys [][]byte
		set     []byte
		scores  []float64
	}{
		keys:    [][]byte{[]byte("location:region1:antibody:SARS-CoV-2:pid:1267449"), []byte("location:region1:antibody:SARS-CoV-2:pid:9784321"), []byte("location:region1:antibody:SARS-CoV-2:pid:2334563")},
		values:  [][]byte{[]byte("positive"), []byte("negative"), []byte("invalid")},
		refKeys: [][]byte{[]byte("false positive"), []byte("refKey2"), []byte("refKey3")},
		set:     []byte("region1"),
		scores:  []float64{2020342, 2021003, 2020350},
	}

	// iterating over example data for verified set of the keys and values
	for i := 0; i < len(antigenData.keys); i++ {
		key1, value1 := antigenData.keys[i], antigenData.values[i]
		tx, err2 := client.VerifiedSet(ctx, key1, value1)
		fmt.Printf("Set and verified key '%s' with value '%s' at tx %d\n", key1, value1, tx.Id)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

	// Zadd the keys and scores to a set called region1
	for i := 0; i < len(antigenData.keys); i++ {
		_, err := client.ZAdd(ctx, antigenData.set, antigenData.scores[i], antigenData.keys[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	// ZScan items by set with a minimal score of 2020343
	var scoremin= &schema.Score{Score: float64(2020343)}
	zStructuredItemList, err := client.ZScan(ctx, &schema.ZScanRequest{
		Set: antigenData.set,
		Limit: uint64(2),
		MinScore: scoremin,
	})
	if err != nil {
		fmt.Println("oupsy")
		log.Fatal(err)
	}

	fmt.Println("ZScan - iterate over a sorted set:")
	for _, item := range zStructuredItemList.GetEntries() {
		fmt.Println(item.String())
		fmt.Println("	------")
	}

	// perform scan by prefix
	prefix := []byte("location:region1")
	scanresult, err := client.Scan(ctx, &schema.ScanRequest{ SeekKey: nil, Prefix: prefix, Limit: 0, SinceTx: 0, Desc: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Scan - print out keys with values (e.g. \"%s\"):\n", prefix)
	fmt.Printf("%v\n", scanresult.GetEntries())
	
	key0Ref := append([]byte("false positive"))
	verifiedIndex, err := client.VerifiedSetReference(ctx, key0Ref, antigenData.keys[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("current state is : %v \n", verifiedIndex)


	// get currentstate
	fmt.Println("Current root - return the last merkle tree root and index stored locally")
	state, err := client.CurrentState(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("current state is : %v \n", state)

	// perform Healthcheck
	err = client.HealthCheck(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
