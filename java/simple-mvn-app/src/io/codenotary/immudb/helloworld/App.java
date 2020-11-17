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

package io.codenotary.immudb.helloworld;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import io.codenotary.immudb4j.FileRootHolder;
import io.codenotary.immudb4j.ImmuClient;
import io.codenotary.immudb4j.KV;
import io.codenotary.immudb4j.KVList;
import io.codenotary.immudb4j.crypto.VerificationException;

public class App {

	public static void main(String[] args) {
		
		ImmuClient client = null;
		
		try {
			
			FileRootHolder rootHolder = FileRootHolder.newBuilder().setRootsFolder("./helloworld_immudb_roots").build();
			
			client = ImmuClient.newBuilder().setServerUrl("localhost").setServerPort(3322).setRootHolder(rootHolder).build();
			
			client.login("immudb", "immudb");

			client.useDatabase("defaultdb");
			
			client.set("hello", "immutable world!".getBytes());
			
			byte[] v = client.safeGet("hello");
			
			System.out.format("(%s, %s)%n", "hello", new String(v));
			

			// Multi-key operations

			KVList.KVListBuilder builder = KVList.newBuilder();

		    builder.add("k123", new byte[]{1, 2, 3});
		    builder.add("k321", new byte[]{3, 2, 1});

		    KVList kvList = builder.build();

		    client.setAll(kvList);


		    List<String> keyList = new ArrayList<>();

		    keyList.add("k123");
		    keyList.add("k321");
		    keyList.add("k231");

		    List<KV> result = client.getAll(keyList);

		    for (KV kv : result) {
		        byte[] key = kv.getKey();
		        byte[] value = kv.getValue();

		        System.out.format("(%s, %s)%n", new String(key), Arrays.toString(value));
		    }

		    // History operations
			client.set("history", new byte[]{1, 2, 3});
			client.set("history", new byte[]{3, 2, 1});

			List<KV> history = client.history("history");

			System.out.format("History of 'history', record 1: (%s, %s)%n", new String(history.get(0).getKey()),
					Arrays.toString(history.get(0).getValue()));
			System.out.format("History of 'history', record 2: (%s, %s)%n", new String(history.get(1).getKey()),
					Arrays.toString(history.get(1).getValue()));

			// Scan operations
            client.set("scanFirstKey", new byte[] {1, 2, 3});
			client.set("scanSecondKey", new byte[] {4, 5, 6});

			List<KV> scan = client.scan("scan", "", 10, false, false);
			System.out.format("Scan results of 'scan', record 1: (%s, %s)%n", new String(scan.get(0).getKey()),
					Arrays.toString(scan.get(0).getValue()));
			System.out.format("Scan results of 'scan', record 2: (%s, %s)%n", new String(scan.get(1).getKey()),
					Arrays.toString(scan.get(1).getValue()));

			// Zadd, zscan operations
			client.zAdd("newSet", "scanFirstKey", 2);
			client.zAdd("newSet", "scanSecondKey", 1);

			List<KV> zScan = client.zScan("newSet", "", 10, false);
			System.out.format("Scan results of 'zscan', record 1: (%s, %s)%n", new String(zScan.get(0).getKey()),
					Arrays.toString(zScan.get(0).getValue()));
			System.out.format("Scan results of 'zscan', record 2: (%s, %s)%n", new String(zScan.get(1).getKey()),
					Arrays.toString(zScan.get(1).getValue()));
		} catch (IOException e) {
			e.printStackTrace();
		} catch (VerificationException e) {
			// Tampering detected!
			// This means the history of changes has been tampered
			e.printStackTrace();
			System.exit(1);
		} finally {
			if (client != null) {
				client.logout();
			}
		}
				
	}
	
}
