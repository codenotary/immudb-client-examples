from immudb import ImmudbClient
import json

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def encode(what: str):
    return what.encode("utf-8")

def decode(what: bytes):
    return what.decode("utf-8")

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    
    # You have to operate on bytes
    setResult = client.set(b'x', b'y')
    print(setResult)            # immudb.datatypes.SetResponse
    print(setResult.id)         # id of transaction
    print(setResult.verified)   # in this case verified = False
                                # see Tamperproof reading and writing

    # Also you get response in bytes
    retrieved = client.get(b'x')
    print(retrieved)        # immudb.datatypes.GetResponse
    print(retrieved.key)    # Value is b'x'
    print(retrieved.value)  # Value is b'y'
    print(retrieved.tx)     # Transaction number

    print(type(retrieved.key))      # <class 'bytes'>
    print(type(retrieved.value))    # <class 'bytes'>

    # Operating with strings
    encodedHello = encode("Hello")
    encodedImmutable = encode("Immutable")
    client.set(encodedHello, encodedImmutable)
    retrieved = client.get(encodedHello)

    print(decode(retrieved.value) == "Immutable")   # Value is True

    notExisting = client.get(b'asdasd')
    print(notExisting)                              # Value is None


    # JSON example
    toSet = {"hello": "immutable"}
    encodedToSet = encode(json.dumps(toSet))
    client.set(encodedHello, encodedToSet)

    retrieved = json.loads(decode(client.get(encodedHello).value))
    print(retrieved)    # Value is {"hello": "immutable"}

    # setAll example - sets all keys to value from dictionary
    toSet = {
        b'1': b'test1',
        b'2': b'test2',
        b'3': b'test3'
    }

    client.setAll(toSet)
    retrieved = client.getAll(list(toSet.keys()))
    print(retrieved) 
    # Value is {b'1': b'test1', b'2': b'test2', b'3': b'test3'}


if __name__ == "__main__":
    main()