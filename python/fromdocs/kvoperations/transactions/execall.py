from immudb import ImmudbClient
from immudb.datatypes import KeyValue, ZAddRequest, ReferenceRequest

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)

    toExecute = [
        KeyValue(b'key', b'value'), 
        ZAddRequest(b'testscore', 100, b'key'),
        KeyValue(b'key2', b'value2'), 
        ZAddRequest(b'testscore', 150, b'key2'),
        ReferenceRequest(b'reference1', b'key')
    ]
    info = client.execAll(toExecute)
    print(info.id) # All in one transaction

    print(client.zScan(b'testscore', b'', 0, 0, True, 10, True, 0, 200)) # Shows these entries
    print(client.get(b'reference1'))

if __name__ == "__main__":
    main()