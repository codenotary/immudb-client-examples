from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)

    state = client.currentState()   # immudb.rootService.State
    print(state.db)         # Current selected DB
    print(state.txId)       # Current transaction ID
    print(state.txHash)     # Current transaction hash
    print(state.signature)  # Current signature

if __name__ == "__main__":
    main()