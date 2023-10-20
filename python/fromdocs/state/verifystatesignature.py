from immudb import ImmudbClient

# All operations are checked against public/private key pair

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)
KEYFILE = "../../../example-public.key"  # Public key path 
                                         # needs immudb server with --signingKey option enabled
                                         # pointing to corresponding private key

def main():
    client = ImmudbClient(URL, publicKeyFile = KEYFILE)
    client.login(LOGIN, PASSWORD, database = DB)
    client.set(b'x', b'1')
    client.verifiedGet(b'x')    # This operation will also fail if public key
                                # is not paired with private one used in immudb

    state = client.currentState()   # immudb.rootService.State
    print(state.db)         # Current selected DB
    print(state.txId)       # Current transaction ID
    print(state.txHash)     # Current transaction hash
    print(state.signature)  # Current signature

if __name__ == "__main__":
    main()