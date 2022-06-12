from grpc import RpcError
from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    first = client.set(b'x', b'y') # Not important, just to demonstrate
                                   # that it will not affect atRevision 
                                   # for other keys

    key = b'immudb130130130'

    client.set(key, b'111')
    client.set(key, b'222')
    client.set(key, b'333')

    print(client.get(key, -2))   # b"111" - value on relative -2 point history
    print(client.get(key, -1))   # b"222" - value on relative -1 point history
    print(client.get(key, 0))    # b"333" - value on relative 0 (current) point history

    print(client.get(key, 1))    # b"111" - value at first revision of key
    print(client.get(key, 2))    # b"222" - value on second revision of key
    print(client.get(key, 3))    # b"333" - value on third revision of key

    try:
        print(client.get(key, -20000))
    except RpcError as error:
        print(error.details()) # invalid key revision number

if __name__ == "__main__":
    main()