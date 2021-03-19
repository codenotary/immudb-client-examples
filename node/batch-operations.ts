import ImmudbClient from 'immudb-node'
import Parameters from 'immudb-node/types/parameters'

const IMMUDB_HOST = '127.0.0.1'
const IMMUDB_PORT = 3322
const IMMUDB_USER = 'immudb'
const IMMUDB_PWD = 'immudb'

const cl = new ImmudbClient({
  host: IMMUDB_HOST,
  port: IMMUDB_PORT,
});

const randNum = Math.floor(Math.random() * Math.floor(100000));
const randStr = `rand${randNum}`;
 
(async () => {
  try {
    // login using the specified username and password
    const loginReq: Parameters.Login = { user: IMMUDB_USER, password: IMMUDB_PWD }
    const loginRes = await cl.login(loginReq)
    console.log('success: login', loginRes)

    // create database
    const createDatabaseReq: Parameters.CreateDatabase = { databasename: randStr }
    const createDatabaseRes = await cl.createDatabase(createDatabaseReq)
    console.log('success: createDatabase', createDatabaseRes)

    // use database just created
    const useDatabaseReq: Parameters.UseDatabase = { databasename: randStr }
    const useDatabaseRes = await cl.useDatabase(useDatabaseReq)
    console.log('success: useDatabase', useDatabaseRes)
    
    // execute a batch insert
    const setAllReq: Parameters.SetAll = { kvsList: [] }
    for (let i = 0; i < 20; i++) {
      setAllReq.kvsList.push({ key: `${i}`, value: `${i}` })
    }
    const setAllRes = await cl.setAll(setAllReq)
    console.log(`success: setAll`, setAllRes)

    // execute a batch read
    const getAllReq: Parameters.GetAll = { keysList: [], sincetx: 0 }
    for (let i = 0; i < 20; i++) {
      getAllReq.keysList.push(`${i}`)
    }
    const getAllRes = await cl.getAll(getAllReq)
    console.log(`success: getAll`, getAllRes)
 
  } catch (err) {
    console.log(err)
  }
})()
