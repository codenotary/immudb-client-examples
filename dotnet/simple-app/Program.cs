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

namespace simple_app;

class Program
{
    public static async Task Main(string[] args)
    {
        FileImmuStateHolder stateHolder = FileImmuStateHolder.NewBuilder()
                .WithStatesFolder("immudb/states")
                .build();

        var client = ImmuClient.Builder()
            .WithStateHolder(stateHolder)
            .WithServerUrl("localhost")
            .WithServerPort(3322)
            .Build();
        await client.Open("immudb", "immudb", "defaultdb");        

        string key = "hello";

        try
        {
            // Setting (adding) a key-value.
            await client.Set(key, Encoding.UTF8.GetBytes("immutable world!"));

            // Getting it back, by key (in a verified way
            // that reports any tampering if it happened).
            Entry entry = await client.VerifiedGet(key);
            Console.WriteLine($"({key}, {Encoding.UTF8.GetString(entry.Value)})\n");
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
        await client.Connection.Pool.Shutdown();
    }
}