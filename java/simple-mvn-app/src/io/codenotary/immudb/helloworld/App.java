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

import io.codenotary.immudb.FileRootHolder;
import io.codenotary.immudb.ImmuClient;
import io.codenotary.immudb.crypto.VerificationException;

public class App {

	public static void main(String args[]) {
		
		ImmuClient client = null;
		
		try {
			
			FileRootHolder rootHolder = FileRootHolder.newBuilder().setCurrentRootsFolder("./helloworld_immudb_roots").build();
			
			client = ImmuClient.newBuilder().setRootHolder(rootHolder).build();
			
			client.login("immudb", "");
			
			
			client.set("hello", "immutable world!".getBytes());
			
			byte[] v = client.safeGet("hello");
			
			System.out.format("(%s, %s)", "hello", new String(v));
			
		} catch (IOException e) {
			e.printStackTrace();
		} catch (VerificationException e) {
			// TODO: tampering detected!
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
