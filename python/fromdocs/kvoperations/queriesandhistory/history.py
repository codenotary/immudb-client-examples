from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    
    client.set(b'test', b'1')
    client.set(b'test', b'2')
    client.set(b'test', b'3')

    history = client.history(b'test', 0, 100, True) # List[immudb.datatypes.historyResponseItem]
    responseItemFirst = history[0]
    print(responseItemFirst.key)    # Entry key (b'test')
    print(responseItemFirst.value)  # Entry value (b'3')
    print(responseItemFirst.tx)     # Transaction id
    
    responseItemThird = history[2]
    print(responseItemThird.key)    # Entry key (b'test')
    print(responseItemThird.value)  # Entry value (b'1')
    print(responseItemThird.tx)     # Transaction id

if __name__ == "__main__":
    main()