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

const ImmudbClient = require('immudb-node')

const IMMUDB_HOST = '127.0.0.1';
const IMMUDB_PORT = '3322';
const IMMUDB_USER = 'immudb';
const IMMUDB_PWD = 'immudb';

const cl = new ImmudbClient.default({ host: IMMUDB_HOST, port: IMMUDB_PORT });

(async () => {
  try {
    const loginReq = { user: IMMUDB_USER, password: IMMUDB_PWD }
    const loginRes = await cl.login(loginReq)
    console.log('success: login:', loginRes)
  
    const useDatabaseReq = { databasename: 'defaultdb' }
    const useDatabaseRes = await cl.useDatabase(useDatabaseReq)
    console.log('success: useDatabase:', useDatabaseRes)
  
    const updateAuthConfigReq = { kind: 1 }
    const updateAuthConfigRes = await cl.updateAuthConfig(updateAuthConfigReq)
    console.log('success: updateAuthConfig', updateAuthConfigRes)
  
    const updateMTLSConfigReq = { enabled: false }
    const updateMTLSConfigRes = await cl.updateMTLSConfig(updateMTLSConfigReq)
    console.log('success: updateMTLSConfig', updateMTLSConfigRes)
  
  } catch (err) {
    console.log(err)
  }
})()
