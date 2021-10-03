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

const IMMUDB_HOST = '127.0.0.1'
const IMMUDB_PORT = 3322
const IMMUDB_USER = 'immudb'
const IMMUDB_PWD = 'immudb'

const cl = new ImmudbClient.default({
  host: IMMUDB_HOST,
  port: IMMUDB_PORT,
  rootPath: 'rootfile'
});

const testDB = 'opsdb';
 
(async () => {
  try {
    // login using the specified username and password
    const loginReq = { user: IMMUDB_USER, password: IMMUDB_PWD }
    const loginRes = await cl.login(loginReq)
    console.log('success: login', loginRes)

    // create database
    const createDatabaseReq = { databasename: testDB }
    const createDatabaseRes = await cl.createDatabase(createDatabaseReq)
    console.log('success: createDatabase', createDatabaseRes)

    // use database just created
    const useDatabaseReq = { databasename: testDB }
    const useDatabaseRes = await cl.useDatabase(useDatabaseReq)
    console.log('success: useDatabase', useDatabaseRes)

    const tableName = `table${ Math.floor(Math.random() * 101) }`;

    // create table using SQLExec
    const sqlExecCreateTableReq = { sql: `create table ${ tableName } (id integer, name varchar, primary key id);` }
    const sqlExecCreateTableRes = await cl.SQLExec(sqlExecCreateTableReq)
    console.log('success: sqlExec create table', sqlExecCreateTableRes)

    // list tables using listTables
    const sqlListTablesRes = await cl.SQLListTables()
    console.log('success: list tables', sqlListTablesRes)
    
    // insert record using SQLExec
    const sqlExecInsertRecordReq = {
      sql: `insert into ${ tableName } (id, name) values (@id, @name);`,
      params: { id: 1, name: 'Joe' }
    }
    const sqlExecInsertRecordRes = await cl.SQLExec(sqlExecInsertRecordReq)
    console.log('success: sqlExec insert record', sqlExecInsertRecordRes)
    
    // query record using SQLQuery
    const sqlQueryReq = {
      sql: `select id,name from ${ tableName } where name=@name;`,
      params: { name: 'Joe' }
    }
    const sqlQueryRes = await cl.SQLQuery(sqlQueryReq)
    console.log('success: sqlQuery query record', sqlQueryRes)
  } catch (err) {
    console.log(err)
  }
})()
