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

package io.codenotary.immudb.helloworld;

import io.codenotary.immudb4j.*;
import io.codenotary.immudb4j.exceptions.CorruptedDataException;
import io.codenotary.immudb4j.exceptions.VerificationException;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

public class App {

    public static void main(String[] args) {

        ImmuClient client = null;

        try {
            FileImmuStateHolder stateHolder = FileImmuStateHolder.newBuilder()
                    .withStatesFolder("./immudb_states")
                    .build();

            client = ImmuClient.newBuilder()
                    .withServerUrl("localhost")
                    .withServerPort(3322)
                    .withStateHolder(stateHolder)
                    .build();

            client.login("immudb", "immudb");

            client.useDatabase("defaultdb");

            String key = "hello";

            try {
                // Setting (adding) a key-value.
                client.set(key, "immutable world!".getBytes());

                // Getting it back, by key (in a verified way
                // that reports any tampering if it happened).
                Entry entry = client.verifiedGet(key);

                byte[] value = entry.getValue();
                System.out.format("(%s, %s)\n", key, new String(value));
            } catch (VerificationException e) {
                // VerificationException means Data Tampering detected!
                // This means the history of changes has been tampered.
                e.printStackTrace();
                System.exit(1);
            }

            String key1 = "key1", key2 = "key2";

            client.set(key1, new byte[]{1, 2, 3});
            client.set(key2, new byte[]{4, 5, 6});

            List<String> keyList = new ArrayList<>();
            keyList.add(key1);
            keyList.add(key2);

            // A multi-key read.
            List<Entry> result = client.getAll(keyList);

            for (Entry e : result) {
                byte[] k = e.getKey();
                byte[] v = e.getValue();

                System.out.format("(%s, %s)\n", new String(k), Arrays.toString(v));
            }

            // History operations.
            client.set("history", new byte[]{1, 2, 3});
            client.set("history", new byte[]{3, 2, 1});

            List<Entry> history = client.history("history", 10, 0, false);

            System.out.format("History of 'history', entry 1: (%s, %s)\n",
                    new String(history.get(0).getKey()),
                    Arrays.toString(history.get(0).getValue())
            );
            System.out.format("History of 'history', entry 2: (%s, %s)\n",
                    new String(history.get(1).getKey()),
                    Arrays.toString(history.get(1).getValue())
            );

            // Scan operations.
            String prefix = "myKey";
            key1 = prefix + "1";
            key2 = prefix + "2";

            client.set(key1, new byte[]{1, 2, 3});

            client.set(key2, new byte[]{4, 5, 6});

            // scan is usually done by a prefix.
            // Of course, we can scan by a complete key name.
            List<Entry> scan = client.scan(prefix, 2, false);

            System.out.format("Scan results of '%s', entry 1: (%s, %s)\n",
                    prefix,
                    new String(scan.get(0).getKey()),
                    Arrays.toString(scan.get(0).getValue())
            );
            System.out.format("Scan results of '%s', entry 2: (%s, %s)\n",
                    prefix,
                    new String(scan.get(1).getKey()),
                    Arrays.toString(scan.get(1).getValue())
            );

            // zAdd, zScan operations.
            String set = "mySet";
            client.zAdd(set, key1, 2);

            client.zAdd(set, key2, 2);

            // Here, we do a zScan providing the `sinceTxId` which should be
            // the latest transaction id we are interested in being considered.
            List<ZEntry> zScan = client.zScan(set, 10, false);

            System.out.format("Results of 'zScan', record 1: (%s, %s)\n",
                    new String(zScan.get(0).getKey()),
                    Arrays.toString(zScan.get(0).getEntry().getValue())
            );
            System.out.format("Results of 'zScan', record 2: (%s, %s)\n",
                    new String(zScan.get(1).getKey()),
                    Arrays.toString(zScan.get(1).getEntry().getValue())
            );

        } catch (IOException | CorruptedDataException e) {
            e.printStackTrace();
        } finally {
            if (client != null) {
                client.logout();
            }
        }

    }

}
