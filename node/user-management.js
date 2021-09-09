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
const { USER_ACTION, USER_PERMISSION } = require('immudb-node/dist/types/user')
const util = require('util')

const IMMUDB_HOST = '127.0.0.1'
const IMMUDB_PORT = 3322
const IMMUDB_USER = 'immudb'
const IMMUDB_PWD = 'immudb'

const cl = new ImmudbClient.default({
  host: IMMUDB_HOST,
  port: IMMUDB_PORT,
})

const randNum = Math.floor(Math.random() * Math.floor(10));
const randStr = `rand${randNum}`;
 
(async () => {
  try {
    const loginReq = { user: IMMUDB_USER, password: IMMUDB_PWD }
    const loginRes = await cl.login(loginReq)
    console.log('success: login', loginRes)

    const createUserReq = {
      user: randStr,
      password: 'Example12#',
      permission: USER_PERMISSION.READ_ONLY,
      database: 'defaultdb',
    }
    const createUserRes = await cl.createUser(createUserReq)
    console.log('success: createUser', createUserRes);

    const listUsersRes = await cl.listUsers()
    console.log('success: listUser', util.inspect(listUsersRes, false, 6, true))

    const changePermissionReq = {
      action: USER_ACTION.GRANT,
      username: randStr,
      database: randStr,
      permission: USER_PERMISSION.READ_WRITE, 
    }
    const changePermissionRes = await cl.changePermission(changePermissionReq)
    console.log('success: changePermission', changePermissionRes);

    const changePasswordReq = {
      user: randStr,
      oldpassword: 'Example12#',
      newpassword: 'Example1234%',
    }
    const changePasswordRes = await cl.changePassword(changePasswordReq)
    console.log('success: changePassword', changePasswordRes);

    const setActiveUserReq = {
      username: randStr,
      active: true,
    }
    const setActiveUserRes = await cl.setActiveUser(setActiveUserReq)
    console.log('success: setActiveUser', setActiveUserRes);

    const logoutRes = await cl.logout()
    console.log('success: logout', logoutRes);

  } catch (err) {
    console.log(err)
  }
})()
