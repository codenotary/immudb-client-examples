from immudb import ImmudbClient
from immudb.client import PersistentRootService

# By default RootService is writing state to RAM
# You can choose different implementation of RootService

# Persistent root service will save to the disk after every verified transaction

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)
PERSISTENT_ROOT_SERVICE_PATH = "/tmp/psr.db" 

def main():
    client = ImmudbClient(URL, rs = PersistentRootService(PERSISTENT_ROOT_SERVICE_PATH))
    client.login(LOGIN, PASSWORD, database = DB)
    client.verifiedSet(b'x', b'1')
    client.verifiedGet(b'x')
    client.verifiedSet(b'x', b'2')
    client.verifiedGet(b'x')

if __name__ == "__main__":
    main()