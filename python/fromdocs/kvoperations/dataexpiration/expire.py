from immudb import ImmudbClient
from datetime import datetime, timedelta
import time

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    client.expireableSet(b"TEST", b"test", datetime.now() + timedelta(seconds=3))
    print(client.get(b"TEST")) # b"test"
    time.sleep(4)
    try:
        print(client.get(b"TEST"))
    except:
        pass # Key not found, because it expires, raises Exception

if __name__ == "__main__":
    main()