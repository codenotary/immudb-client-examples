from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    client.verifiedSet(b'x', b'1') 
    client.verifiedSet(b'y', b'1') 
    retrieved = client.verifiedGet(b'x') 
    print(retrieved.refkey)     # Entry reference key (None)

    client.verifiedSetReference(b'x', b'reference1')
    client.setReference(b'x', b'reference2')
    client.setReference(b'y', b'reference2')
    client.verifiedSet(b'y', b'2') 

    retrieved = client.verifiedGet(b'reference1')
    print(retrieved.key)        # Entry key (b'x')
    print(retrieved.refkey)     # Entry reference key (b'reference1')
    print(retrieved.verified)   # Entry verification status (True)

    retrieved = client.verifiedGet(b'reference2')
    print(retrieved.key)        # Entry key (b'y')
    print(retrieved.refkey)     # Entry reference key (b'reference2')
    print(retrieved.verified)   # Entry verification status (True)
    print(retrieved.value)      # Entry value (b'3')

    retrieved = client.verifiedGet(b'x')
    print(retrieved.key)        # Entry key (b'x')
    print(retrieved.refkey)     # Entry reference key (None)
    print(retrieved.verified)   # Entry verification status (True)

    retrieved = client.get(b'reference2')
    print(retrieved.key)        # Entry key (b'y')

if __name__ == "__main__":
    main()