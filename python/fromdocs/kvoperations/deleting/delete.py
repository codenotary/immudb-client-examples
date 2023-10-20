from immudb import ImmudbClient
from immudb.datatypes import DeleteKeysRequest

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    client.set(b"immu", b"immudb-not-rulezz")
    print(client.get(b"immu"))  # b"immudb-not-rulezz"

    deleteRequest = DeleteKeysRequest(keys = [b"immu"])
    client.delete(deleteRequest)
    print(client.get(b"immu"))  # None

if __name__ == "__main__":
    main()