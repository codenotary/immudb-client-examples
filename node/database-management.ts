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

import ImmudbClient from 'immudb-node'
import Parameters from 'immudb-node/dist/types/parameters'
import util from 'util'

const IMMUDB_HOST = '127.0.0.1'
const IMMUDB_PORT = 3322
const IMMUDB_USER = 'immudb'
const IMMUDB_PWD = 'immudb'

const cl = new ImmudbClient({
  host: IMMUDB_HOST,
  port: IMMUDB_PORT,
});

const randNum = Math.floor(Math.random() * Math.floor(10));
const randStr = `rand${randNum}`;
 
(async () => {
  try {
    const loginReq: Parameters.Login = { user: IMMUDB_USER, password: IMMUDB_PWD }
    const loginRes = await cl.login(loginReq)
    console.log('success: login', loginRes);

    const createDatabaseReq: Parameters.CreateDatabase = { databasename: randStr }
    const createDatabaseRes = await cl.createDatabase(createDatabaseReq)
    console.log('success: createDatabase', createDatabaseRes);

    const useDatabaseReq: Parameters.UseDatabase = { databasename: randStr }
    const useDatabaseRes = await cl.useDatabase(useDatabaseReq)
    console.log('success: useDatabase', useDatabaseRes);

    const setReq: Parameters.Set = { key: randStr, value: randStr }
    const setRes = await cl.set(setReq)
    console.log('success: set', setRes);

    const listDatabasesRes = await cl.listDatabases()
    console.log('success: listDatabases', listDatabasesRes);

    console.log(util.inspect(listDatabasesRes, false, 8, true))

    const healthRes = await cl.health()
    console.log('success: health', healthRes);

  } catch (err) {
    console.log(err)
  }
})()
