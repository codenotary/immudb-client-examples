/*
Copyright 2022 CodeNotary, Inc. All rights reserved.

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

using System.Text;
using ImmuDB;
using ImmuDB.Exceptions;
using ImmuDB.SQL;

namespace simple_app;

class Program
{
    private static string immudbServerAddress = "localhost";
    
    public static async Task Main(string[] args)
    {
        string? envAddr = Environment.GetEnvironmentVariable("IMMUDB_ADDRESS");
        if(!string.IsNullOrEmpty(envAddr)) {
            immudbServerAddress = envAddr;
        }
        await OpenConnectionExample();
        await AnotherOpenConnectionExample();
        await GetSetScanUsageExample();
        await SetAllGetAllExample();
        await SqlUsageExample();

        SyncOpenConnectionExample();
        SyncSqlUsageExample();

        await ImmuClient.ReleaseSdkResources();
    }

    private static async Task OpenConnectionExample()
    {
        var client = ImmuClient.NewBuilder().WithServerUrl(immudbServerAddress).Build();
        await client.Open("immudb", "immudb", "defaultdb");

        string key = "hello";

        try
        {
            await client.VerifiedSet(key, "immutable world!");

            // Getting it back, by key (in a verified way that reports any tampering if it happened).
            Entry entry = await client.VerifiedGet(key);
            Console.WriteLine($"{key}, {entry.ToString()}");
        }
        catch (VerificationException e)
        {
            // VerificationException means Data Tampering detected!
            // This means the history of changes has been tampered.
            Console.WriteLine(e.ToString());
        }
        await client.Close();
    }
    
    private static void SyncOpenConnectionExample()
    {
        var client = ImmuClientSync.NewBuilder().WithServerUrl(immudbServerAddress).Build();
        client.Open("immudb", "immudb", "defaultdb");

        string key = "hello";

        try
        {
            client.VerifiedSet(key, "immutable world!");

            // Getting it back, by key (in a verified way that reports any tampering if it happened).
            Entry entry = client.VerifiedGet(key);
            Console.WriteLine($"{key}, {entry.ToString()}");
        }
        catch (VerificationException e)
        {
            // VerificationException means Data Tampering detected!
            // This means the history of changes has been tampered.
            Console.WriteLine(e.ToString());
        }
        client.Close();
    }

    private static async Task AnotherOpenConnectionExample()
    {
        var client = await ImmuClient.NewBuilder().WithServerUrl(immudbServerAddress).Open();
        string key = "hello";

        try
        {
            // Setting (adding) a key-value.
            await client.VerifiedSet(key, "immutable world!");

            // Getting it back, by key (in a verified way that reports any tampering if it happened).
            Entry entry = await client.VerifiedGet(key);
            Console.WriteLine($"{key}, {entry.ToString()}");
        }
        catch (VerificationException e)
        {
            // VerificationException means Data Tampering detected!
            // This means the history of changes has been tampered.
            Console.WriteLine(e.ToString());
        }

        await client.Close();
    }

    private static async Task SetAllGetAllExample()
    {
        var client = new ImmuClient(immudbServerAddress, 3322);
        await client.Open("immudb", "immudb", "defaultdb");


        try
        {
            List<string> keys = new List<string>();
            keys.Add("k0");
            keys.Add("k1");

            List<byte[]> values = new List<byte[]>();
            values.Add(new byte[] { 0, 1, 0, 1 });
            values.Add(new byte[] { 1, 0, 1, 0 });

            List<KVPair> kvListBuilder = new List<KVPair>();

            for (int i = 0; i < keys.Count; i++)
            {
                kvListBuilder.Add(new KVPair(keys[i], values[i]));
            }
            await client.SetAll(kvListBuilder);

            List<Entry> getAllResult = await client.GetAll(keys);

            for (int i = 0; i < getAllResult.Count; i++)
            {
                Entry entry = getAllResult[i];
                Console.WriteLine($"({string.Join(" ", entry.Key)}, {keys[i]}):({string.Join(" ", entry.Value)}, {string.Join(" ", values[i])})");
            }

            await client.Close();
        }
        catch (VerificationException e)
        {
            // VerificationException means Data Tampering detected!
            // This means the history of changes has been tampered.
            Console.WriteLine(e.ToString());
        }

        await client.Close();
    }


    private static async Task GetSetScanUsageExample()
    {
        var client = ImmuClient.NewBuilder()
            .WithServerUrl(immudbServerAddress)
            .WithServerPort(3322)
            .Build();
        await client.Open("immudb", "immudb", "defaultdb");

        string key = "hello";

        try
        {
            // Setting (adding) a key-value.
            await client.Set(key, "immutable world!");
            var immuWorldVal = await client.Get(key);
            Console.WriteLine($"{key} : {immuWorldVal.ToString()}");

            // Getting it back, by key (in a verified way that reports any tampering if it happened).
            Entry entry = await client.VerifiedGet(key);
            Console.WriteLine($"({key}, {entry.ToString()})\n");
        }
        catch (VerificationException e)
        {
            // VerificationException means Data Tampering detected!
            // This means the history of changes has been tampered.
            Console.WriteLine(e.ToString());
            Environment.Exit(1);
        }

        string key1 = "key1", key2 = "key2";

        await client.Set(key1, new byte[] { 1, 2, 3 });
        await client.Set(key2, new byte[] { 4, 5, 6 });

        List<string> keyList = new List<string>();
        keyList.Add(key1);
        keyList.Add(key2);

        // A multi-key read.
        List<Entry> result = await client.GetAll(keyList);

        foreach (Entry e in result)
        {
            Console.WriteLine($"({string.Join(" ", key)}, {string.Join(" ", e.Value)})\n");
        }

        // History operations.
        await client.Set("history", new byte[] { 1, 2, 3 });
        await client.Set("history", new byte[] { 3, 2, 1 });

        List<Entry> history = await client.History("history", 10, 0, false);

        Console.WriteLine($"History of 'history', entry 1: ({Encoding.UTF8.GetString(history[0].Key)}, {string.Join(" ", history[0].Value)})\n");
        Console.WriteLine($"History of 'history', entry 2: ({Encoding.UTF8.GetString(history[1].Key)}, {string.Join(" ", history[1].Value)})\n");

        // Scan operations.
        String prefix = "myKey";
        key1 = prefix + "1";
        key2 = prefix + "2";

        await client.Set(key1, new byte[] { 1, 2, 3 });

        await client.Set(key2, new byte[] { 4, 5, 6 });

        // scan is usually done by a prefix.
        // Of course, we can scan by a complete key name.
        List<Entry> scan = await client.Scan(prefix, 2, false);

        Console.WriteLine($"Scan results of '{prefix}', entry 1: ({Encoding.UTF8.GetString(scan[0].Key)}, {string.Join(" ", scan[0].Value)})\n");
        Console.WriteLine($"Scan results of '{prefix}', entry 2: ({Encoding.UTF8.GetString(scan[1].Key)}, {string.Join(" ", scan[1].Value)})\n");

        // zAdd, zScan operations.
        String set = "mySet";
        await client.ZAdd(set, key1, 2);
        await client.ZAdd(set, key2, 2);

        // Here, we do a zScan providing the `sinceTxId` which should be
        // the latest transaction id we are interested in being considered.
        List<ZEntry> zScan = await client.ZScan(set, 10, false);

        Console.WriteLine($"Results of 'zScan', record 1: ({Encoding.UTF8.GetString(zScan[0].Key)},{string.Join(" ", zScan[0].Entry.Value)})\n");
        Console.WriteLine($"Results of 'zScan', record 2: ({Encoding.UTF8.GetString(zScan[1].Key)},{string.Join(" ", zScan[1].Entry.Value)})\n");

        await client.Close();
    }

    private static async Task SqlUsageExample()
    {
        var client = new ImmuClient(immudbServerAddress, 3322);
        await client.Open("immudb", "immudb", "defaultdb");

        await client.SQLExec("CREATE TABLE IF NOT EXISTS logs(id INTEGER AUTO_INCREMENT, created TIMESTAMP, entry VARCHAR, PRIMARY KEY id)");
        await client.SQLExec("CREATE INDEX IF NOT EXISTS ON logs(created)");
        var rspInsert = await client.SQLExec("INSERT INTO logs(created, entry) VALUES($1, $2)",
                SQLParameter.Create(DateTime.UtcNow),
                SQLParameter.Create("hello immutable world"));
        var queryResult = await client.SQLQuery("SELECT created, entry FROM LOGS order by created DESC");
        var sqlVal = queryResult.Rows[0]["entry"];
        
        Console.WriteLine($"The log entry is: {sqlVal.Value.ToString()}");
        await client.Close();
    }
    
    private static void SyncSqlUsageExample()
    {
        var client = new ImmuClientSync(immudbServerAddress, 3322);
        client.Open("immudb", "immudb", "defaultdb");        

        client.SQLExec("CREATE TABLE IF NOT EXISTS logs(id INTEGER AUTO_INCREMENT, created TIMESTAMP, entry VARCHAR, PRIMARY KEY id)");
        client.SQLExec("CREATE INDEX IF NOT EXISTS ON logs(created)");
        var rspInsert = client.SQLExec("INSERT INTO logs(created, entry) VALUES($1, $2)",
                SQLParameter.Create(DateTime.UtcNow),
                SQLParameter.Create("hello immutable world"));
        var queryResult = client.SQLQuery("SELECT created, entry FROM LOGS order by created DESC");
        var sqlVal = queryResult.Rows[0]["entry"];
        
        Console.WriteLine($"The log entry is: {sqlVal.Value.ToString()}");
        client.Close();
    }

    
}