from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    toSet = {
        b"aaa": b'1',
        b'bbb': b'2',
        b'ccc': b'3',
        b'acc': b'1',
        b'aac': b'2',
        b'aac:test1': b'3',
        b'aac:test2': b'1',
        b'aac:xxx:test': b'2'
    }
    client.setAll(toSet)
    
    result = client.scan(b'', b'', True, 100) # All entries
    print(result)
    result = client.scan(b'', b'aac', True, 100) # All entries with prefix 'aac' including 'aac'
    print(result)

    # Seek key example (allows retrieve entries in proper chunks):
    result = client.scan(b'', b'', False, 3)
    while result:
        for item, value in result.items():
            print("SEEK", item, value)
        lastKey = list(result.keys())[-1]
        result = client.scan(lastKey, b'', False, 3)

if __name__ == "__main__":
    main()