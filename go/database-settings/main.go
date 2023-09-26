/*
Copyright 2022 Codenotary Inc. All rights reserved.

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
	"flag"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
)

// go mod tidy
// go build
// ./db_settings -h
// e.g. ./db_settings -db=mydb -create=true

var config struct {
	Addr            string
	Port            int
	Username        string
	Password        string
	DBName          string
	CreateDatabase  bool
	Replica         bool
	PrimaryUser     string
	PrimaryPassword string
	PrimaryDB       string
	PrimaryHost     string
	PrimaryPort     int
}

func init() {
	flag.StringVar(&config.Addr, "addr", "127.0.0.1", "IP address of immudb server")
	flag.IntVar(&config.Port, "port", 3322, "Port number of immudb server")
	flag.StringVar(&config.Username, "user", "immudb", "Username for authenticating to immudb")
	flag.StringVar(&config.Password, "pass", "immudb", "Password for authenticating to immudb")
	flag.StringVar(&config.DBName, "db", "defaultdb", "Name of the database to use")
	flag.BoolVar(&config.CreateDatabase, "create", true, "Create database")

	flag.BoolVar(&config.Replica, "replica", false, "Create a sync replica")
	flag.StringVar(&config.PrimaryUser, "primary-user", "immudb", "Username for authenticating to primary immudb")
	flag.StringVar(&config.PrimaryPassword, "primary-pass", "immudb", "Password for authenticating to primary immudb")
	flag.StringVar(&config.PrimaryDB, "primary-db", "primary-db", "Name of the primary database to replicate")
	flag.StringVar(&config.PrimaryHost, "primary-addr", "127.0.0.1", "IP address of primary immudb server")
	flag.IntVar(&config.PrimaryPort, "primary-port", 3322, "Port number of immudb server")
	flag.Parse()
}

func main() {
	// even though the server address and port are defaults, setting them as a reference
	opts := immudb.DefaultOptions().WithAddress(config.Addr).WithPort(config.Port)

	client := immudb.NewClient().WithOptions(opts)

	err := client.OpenSession(context.Background(), []byte(config.Username), []byte(config.Password), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	defer client.CloseSession(context.Background())
	var repl schema.ReplicationNullableSettings = schema.ReplicationNullableSettings{
		Replica: &schema.NullableBool{Value: false},
	}
	if config.Replica {
		repl.Replica.Value = true
		repl.SyncReplication = &schema.NullableBool{Value: true}
		repl.PrimaryDatabase = &schema.NullableString{Value: config.PrimaryDB}
		repl.PrimaryHost = &schema.NullableString{Value: config.PrimaryHost}
		repl.PrimaryPort = &schema.NullableUint32{Value: uint32(config.PrimaryPort)}
		repl.PrimaryUsername = &schema.NullableString{Value: config.PrimaryUser}
		repl.PrimaryPassword = &schema.NullableString{Value: config.PrimaryPassword}
	}
	dbSettings := &schema.DatabaseNullableSettings{
		ReplicationSettings:     &repl,
		ExcludeCommitTime:       &schema.NullableBool{Value: false},
		MaxConcurrency:          &schema.NullableUint32{Value: 1_000},
		MaxIOConcurrency:        &schema.NullableUint32{Value: 1},
		TxLogCacheSize:          &schema.NullableUint32{Value: 100_000},
		VLogMaxOpenedFiles:      &schema.NullableUint32{Value: 30},
		TxLogMaxOpenedFiles:     &schema.NullableUint32{Value: 30},
		CommitLogMaxOpenedFiles: &schema.NullableUint32{Value: 5},
		IndexSettings: &schema.IndexNullableSettings{
			FlushThreshold:           &schema.NullableUint32{Value: 5_000_000},
			SyncThreshold:            &schema.NullableUint32{Value: 10_000_000},
			CacheSize:                &schema.NullableUint32{Value: 1_000_000},
			MaxActiveSnapshots:       &schema.NullableUint32{Value: 100},
			RenewSnapRootAfter:       &schema.NullableUint64{Value: 0},
			CompactionThld:           &schema.NullableUint32{Value: 1_000_000},
			DelayDuringCompaction:    &schema.NullableUint32{Value: 10},
			NodesLogMaxOpenedFiles:   &schema.NullableUint32{Value: 30},
			HistoryLogMaxOpenedFiles: &schema.NullableUint32{Value: 15},
			CommitLogMaxOpenedFiles:  &schema.NullableUint32{Value: 5},
			FlushBufferSize:          &schema.NullableUint32{Value: 4096},
			CleanupPercentage:        &schema.NullableFloat{Value: 1},
		},
		ReadTxPoolSize:  &schema.NullableUint32{Value: 1_000},
		SyncFrequency:   &schema.NullableMilliseconds{Value: 20},
		WriteBufferSize: &schema.NullableUint32{Value: 16_777_216},
		AhtSettings: &schema.AHTNullableSettings{
			SyncThreshold:   &schema.NullableUint32{Value: 10_000_000},
			WriteBufferSize: &schema.NullableUint32{Value: 4096},
		},
		MaxActiveTransactions: &schema.NullableUint32{Value: 10_000},
	}

	if config.CreateDatabase {
		// below settings can only be set at database creation
		dbSettings.FileSize = &schema.NullableUint32{Value: 1 << 30} //1024Mb
		dbSettings.MaxKeyLen = &schema.NullableUint32{Value: 40}
		dbSettings.MaxValueLen = &schema.NullableUint32{Value: 256}
		dbSettings.MaxTxEntries = &schema.NullableUint32{Value: 1024}

		dbSettings.IndexSettings.MaxNodeSize = &schema.NullableUint32{Value: 16384} // 16Kb

		_, err = client.CreateDatabaseV2(context.Background(), config.DBName, dbSettings)
		fmt.Printf("Sucessfully created database: '%s'\n", config.DBName)

	} else {
		_, err = client.UpdateDatabaseV2(context.Background(), config.DBName, dbSettings)
		fmt.Printf("Sucessfully updated database: '%s'\n", config.DBName)
	}
	if err != nil {
		log.Fatal(err)
	}

}
