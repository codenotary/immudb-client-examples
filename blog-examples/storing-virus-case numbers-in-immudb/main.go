package main

import (
	"context"
	"fmt"
	immuapi "github.com/codenotary/immudb/pkg/api"
	"github.com/codenotary/immudb/pkg/api/schema"
	immuschema "github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func printItem(key []byte, value []byte, message interface{}) {
	var index uint64
	ts := uint64(time.Now().Unix())
	var verified, isVerified bool
	var hash []byte
	switch m := message.(type) {
	case *immuschema.Index:
		index = m.Index
		dig := immuapi.Digest(index, key, value)
		hash = dig[:]
	case *immuclient.VerifiedIndex:
		index = m.Index
		dig := immuapi.Digest(index, key, value)
		hash = dig[:]
		verified = m.Verified
		isVerified = true
	case *immuschema.Item:
		key = m.Key
		value = m.Value
		index = m.Index
		hash = m.Hash()
	case *immuschema.StructuredItem:
		key = m.Key
		value = m.Value.Payload
		ts = m.Value.Timestamp
		index = m.Index
		hash, _ = m.Hash()
	case *immuschema.ZStructuredItem:
		key = m.Item.Key
		value = m.Item.Value.Payload
		ts = m.Item.Value.Timestamp
		index = m.Item.Index
		hash, _ = m.Item.Hash()
	case *immuclient.VerifiedItem:
		key = m.Key
		value = m.Value
		index = m.Index
		ts = m.Time
		verified = m.Verified
		isVerified = true
		me, _ := immuschema.Merge(value, ts)
		dig := immuapi.Digest(index, key, me)
		hash = dig[:]
	}
	if !isVerified {
		fmt.Printf("index: %d\n key: %s\n value: %s\n hash: %x\n time: %s\n",
			index,
			key,
			value,
			hash,
			time.Unix(int64(ts), 0))
		return
	}
	fmt.Printf("index: %d\n key: %s\n value: %s\n hash: %x\n time: %s\n verified: %t\n",
		index,
		key,
		value,
		hash,
		time.Unix(int64(ts), 0),
		verified)
}


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

	for i := 0; i < len(antigenData.keys); i++ {
		key1, value1 := antigenData.keys[i], antigenData.values[i]
		_, err2 := client.Set(ctx, key1, value1)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

	//len(antigenData.scores)
	for i := 0; i < len(antigenData.keys); i++ {
		_, err := client.ZAdd(ctx, antigenData.set, antigenData.scores[i], antigenData.keys[i], nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	/* How to use ZScanOptions
	type ZScanOptions struct {
		Set                  []byte   `protobuf:"bytes,1,opt,name=set,proto3" json:"set,omitempty"`
		Offset               []byte   `protobuf:"bytes,2,opt,name=offset,proto3" json:"offset,omitempty"`
		Limit                uint64   `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
		Reverse              bool     `protobuf:"varint,4,opt,name=reverse,proto3" json:"reverse,omitempty"`
		Min                  *Score   `protobuf:"bytes,5,opt,name=min,proto3" json:"min,omitempty"`
		Max                  *Score   `protobuf:"bytes,6,opt,name=max,proto3" json:"max,omitempty"`
		XXX_NoUnkeyedLiteral struct{} `json:"-"`
		XXX_unrecognized     []byte   `json:"-"`
		XXX_sizecache        int32    `json:"-"`
	}
	*/

	var scoremin= &schema.Score{Score: float64(2020343)}
	zStructuredItemList, err := client.ZScan(ctx, &schema.ZScanOptions{
		Set: antigenData.set,
		Limit: uint64(2),
		Min: scoremin,
	})
	if err != nil {
		fmt.Println("oupsy")
		log.Fatal(err)
	}
	fmt.Println("ZScan - iterate over a sorted set:")
	for _, item := range zStructuredItemList.Items {
		printItem(nil, nil, item)
		fmt.Println("	------")
	}

	//------> Scan
	prefix := []byte("location:region1")

	scanresult, err := client.Scan(ctx, &schema.ScanOptions{ Prefix: []byte("location:region1"), Limit: uint64(3), Reverse: true})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Scan - iterate over keys having the specified prefix (e.g. \"%s\"):\n", prefix)
	for _, item := range scanresult.Items {
		printItem(nil, nil, item)
		fmt.Println("	------")
	}


	fmt.Println("Current root - return the last merkle tree root and index stored locally")
	currentRoot, err := client.CurrentRoot(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if currentRoot == nil {
		fmt.Println("no root found: immudb is empty")
	}
	fmt.Printf("index: %s\n hash:  %s\n", currentRoot.Payload, currentRoot.Signature)


	//------> SafeReference
	key0Ref := append([]byte("false positive"))
	verifiedIndex, err := client.Reference(ctx, key0Ref, antigenData.keys[0],nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SafeReference - add and verify a reference key to an existing entry:")
	printItem(antigenData.keys[0], antigenData.values[0], verifiedIndex)

}
