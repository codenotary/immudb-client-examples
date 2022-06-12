from grpc import RpcError
from immudb import ImmudbClient
from immudb.constants import PERMISSION_ADMIN, PERMISSION_R, PERMISSION_RW
from immudb.grpc.schema_pb2 import GRANT, REVOKE
from enum import IntEnum

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)
    passwordForNewUsers = "Te1st!@#Test"
    try:
        client.createUser("tester1", passwordForNewUsers, PERMISSION_R, DB)
        client.createUser("tester2", passwordForNewUsers, PERMISSION_RW, DB)
        client.createUser("tester3", passwordForNewUsers, PERMISSION_ADMIN, DB)
    except RpcError as exception:
        print(exception.details())

    users = client.listUsers().userlist.users # immudb.handler.listUsers.listUsersResponse
    for user in users:
        print("User", user.user)
        print("Created by", user.createdby)
        print("Creation date", user.createdat)
        print("Is active", user.active)
        for permission in user.permissions:
            print("Permission", permission.database, permission.permission)
        print("---")

    client.login("tester3", passwordForNewUsers, DB)
    client.changePermission(GRANT, "tester2", DB, PERMISSION_ADMIN)
    client.changePermission(REVOKE, "tester2", DB, PERMISSION_ADMIN)

    client.login(LOGIN, PASSWORD, database = DB)
    # Changing password
    client.changePassword("tester1", "N1ewpassword!", passwordForNewUsers)

    # User logs with new password
    client.login("tester1", "N1ewpassword!")

    client.login(LOGIN, PASSWORD, database = DB)
    client.changePassword("tester1", passwordForNewUsers, "N1ewpassword!")
    

    client.login("tester1", passwordForNewUsers, DB)

    # No permissions to write
    try:
        client.set(b"test", b"test")
    except RpcError as exception:
        print(exception.details())

    # But has permissions to read
    result = client.get(b"test")

    client.login("tester3", passwordForNewUsers, DB)

    # Now will have permissions to write
    client.changePermission(GRANT, "tester1", DB, PERMISSION_RW)
    client.login("tester1", passwordForNewUsers, DB)
    client.set(b"test", b"test")
    result = client.get(b"test")

    client.login("tester3", passwordForNewUsers, DB)

    # Now will have permissions to nothing
    client.changePermission(REVOKE, "tester1", DB, PERMISSION_RW)

    try:
        client.login("tester1", passwordForNewUsers, DB)
    except RpcError as exception:
        print(exception.details())
    
    client.login("tester3", passwordForNewUsers, DB)
    client.changePermission(GRANT, "tester1", DB, PERMISSION_RW)


if __name__ == "__main__":
    main()