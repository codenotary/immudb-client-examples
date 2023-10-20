from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    dictToSetGet = {
        b'key1': b'value1',
        b'key2': b'value2',
        b'key3': b'value3'
    }
    response = client.setAll(dictToSetGet)
    print(response.id) # All in one transaction

    response = client.getAll([b'key1', b'key2', b'key3'])
    print(response) # The same as dictToSetGet, retrieved in one step

if __name__ == "__main__":
    main()