
from grpc import RpcError
from immudb import ImmudbClient

URL = "localhost:3322"  # immudb running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)


def main():
    client = ImmudbClient(URL)
    # database parameter is optional
    client.login(LOGIN, PASSWORD, database=DB)
    client.logout()

    # Bad login
    try:
        client.login("verybadlogin", "verybadpassword")
    except RpcError as exception:
        print(exception.debug_error_string())
        print(exception.details())


    # Managed session support
    with client.openManagedSession(LOGIN, PASSWORD, database=DB) as session:
        transaction = session.newTx()
        transaction.sqlExec("CREATE TABLE IF NOT EXISTS connectiontest1 (id INTEGER AUTO_INCREMENT, name VARCHAR[255], PRIMARY KEY(id))")
        commited = transaction.commit()
        print(commited) # If table was created exists - shows transaction informations

    # Not managed session. You need to handle keep alive request yourself
    session = client.openSession(LOGIN, PASSWORD, database=DB)
    transaction = session.newTx()
    transaction.sqlExec("CREATE TABLE IF NOT EXISTS connectiontest2 (id INTEGER AUTO_INCREMENT, name VARCHAR[255], PRIMARY KEY(id))")
    commited = transaction.commit()
    print(commited) # If table was created - shows transaction informations
    client.closeSession()



if __name__ == "__main__":
    main()
