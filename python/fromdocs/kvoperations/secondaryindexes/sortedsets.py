from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    client.set(b"user1", b"user1@mail.com")
    client.set(b"user2", b"user2@mail.com")
    client.set(b"user3", b"user3@mail.com")
    client.set(b"user4", b"user3@mail.com")

    client.zAdd(b"age", 100, b"user1")
    client.zAdd(b"age", 101, b"user2")
    client.zAdd(b"age", 99, b"user3")
    client.zAdd(b"age", 100, b"user4")

    scanResult = client.zScan(b"age", b"", 0, 0, True, 50, False, 100, 101)
    print(scanResult)   # Shows records with 'age' 100 <= score < 101
                        # with descending order and limit = 50


if __name__ == "__main__":
    main()