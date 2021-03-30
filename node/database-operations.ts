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

const IMMUDB_HOST = '127.0.0.1'
const IMMUDB_PORT = 3322
const IMMUDB_USER = 'immudb'
const IMMUDB_PWD = 'immudb'

const cl = new ImmudbClient({
  host: IMMUDB_HOST,
  port: IMMUDB_PORT,
  rootPath: 'rootfile'
});

const rand = '1';
const testDB = 'opsdb';
 
(async () => {
  try {
    // login using the specified username and password
    const loginReq: Parameters.Login = { user: IMMUDB_USER, password: IMMUDB_PWD }
    const loginRes = await cl.login(loginReq)
    console.log('success: login', loginRes)

    // create database
    const createDatabaseReq: Parameters.CreateDatabase = { databasename: testDB }
    const createDatabaseRes = await cl.createDatabase(createDatabaseReq)
    console.log('success: createDatabase', createDatabaseRes)

    // use database just created
    const useDatabaseReq: Parameters.UseDatabase = { databasename: testDB }
    const useDatabaseRes = await cl.useDatabase(useDatabaseReq)
    console.log('success: useDatabase', useDatabaseRes)

    // add new item having the specified key and value
    const setReq: Parameters.Set = { key: rand, value: rand }
    const setRes = await cl.set(setReq)
    console.log('success: set', setRes)

    const index = setRes?.id
    if (!index) {
      throw new Error()
    }

    // get item having the specified key
    const getReq: Parameters.Get = { key: rand }
    const getRes = await cl.get(getReq)
    console.log('success: get', getRes)

    // increase occurences of items having the
    // same key
    for (let i = 0; i < 10; i++) {
      const setReq: Parameters.Set = { key: rand, value: rand }
      await cl.set(setReq)
    }

    // iterate over keys having the specified
    // prefix
    const scanReq: Parameters.Scan = {
      seekkey: rand,
      prefix: rand,
      desc: false,
      limit: 1,
      sincetx: 0,
      nowait: false,
    }
    const scanRes = await cl.scan(scanReq)
    console.log('success: scan', scanRes)

    // return an element by index
    const txByIdReq: Parameters.TxById = { tx: index }
    const txByIdRes = await cl.txById(txByIdReq)
    console.log('success: txById', txByIdRes)

    // fetch paginated history for the item having the
    // specified key
    const historyReq: Parameters.History = {
      key: rand,
      offset: 10,
      limit: 5,
      desc: false,
      sincetx: 0
    }
    const historyRes = await cl.history(historyReq)
    console.log('success: history', historyRes)

    // iterate over a sorted set
    const zScanReq: Parameters.ZScan = {
      set: rand,
      seekkey: rand,
      seekscore: 10,
      seekattx: index,
      inclusiveseek: false,
      limit: 5,
      desc: false,
      sincetx: 0,
      nowait: true
    }
    const zScanRes = await cl.zScan(zScanReq)
    console.log('success: zScan', zScanRes)

    // execute a batch read
    const getAllReq: Parameters.GetAll = {
      keysList: [rand],
      sincetx: 0
    }
    const getAllRes = await cl.getAll(getAllReq)
    console.log('success: getBatch', getAllRes)

    // check immudb health status
    const healthRes = await cl.health()
    console.log('success: health', healthRes)

    // get current root info
    const currentStateRes = await cl.currentState()
    console.log('success: currentState', currentStateRes)

    // safely add new item having the specified key and value
    const verifiedSetReq: Parameters.VerifiedSet = {
      key: rand+10,
      value: rand+10,
    }
    const verifiedSetRes = await cl.verifiedSet(verifiedSetReq)
    console.log('success: verifiedSet', verifiedSetRes)

    // get current root info
    const currentStateRes2 = await cl.currentState()
    console.log('success: currentState', currentStateRes2)

    // safely add new item having the specified key and value
    const verifiedSetReq2: Parameters.VerifiedSet = {
      key: rand+11,
      value: rand+11,
    }
    const verifiedSetRes2 = await cl.verifiedSet(verifiedSetReq2)
    console.log('success: verifiedSet', verifiedSetRes2)

    // safely add new item having the specified key and value
    const verifiedSetReq3: Parameters.VerifiedSet = {
      key: rand+12,
      value: rand+12,
    }
    const verifiedSetRes3 = await cl.verifiedSet(verifiedSetReq3)
    console.log('success: verifiedSet', verifiedSetRes3)

    // safely get item by key
    const verifiedGetReq: Parameters.VerifiedGet = {
      key: rand+12,
    }
    const verifiedGetRes = await cl.verifiedGet(verifiedGetReq)
    console.log('success: verifiedGet', verifiedGetRes)

  } catch (err) {
    console.log(err)
  }
})()
