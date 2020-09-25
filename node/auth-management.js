const ImmudbClient = require('immudb-node')

ImmudbClient({
  address: '127.0.0.1:3322',
}, main)

const rand = '' + Math.floor(Math.random()
  * Math.floor(100000))
 
async function main(err, cl) {
  if (err) {
    return console.log(err)
  }

  try {
    let req = { username: 'immudb', password: 'immudb' }
    let res = await cl.login(req)

    res = await cl.useDatabase({ database: 'defaultdb' })

    await cl.updateAuthConfig({ auth: types.auth.disabled })

    await cl.updateMTLSConfig({ enabled: false })

  } catch (err) {
    console.log(err)
  }
}
