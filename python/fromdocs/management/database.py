from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    testDatabase = "test"
    
    databases = client.databaseList()
    if(testDatabase not in databases):
        client.createDatabase(testDatabase)

    client.useDatabase("test")
    client.set(b"test", b"test")
    client.databaseList()
    print(client.get(b"test"))

if __name__ == "__main__":
    main()