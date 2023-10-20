from immudb import ImmudbClient
from uuid import uuid4

URL = "localhost:3322"  # ImmuDB running on your machine
LOGIN = "immudb"        # Default username
PASSWORD = "immudb"     # Default password
DB = b"defaultdb"       # Default database name (must be in bytes)

def main():
    client = ImmudbClient(URL)
    client.login(LOGIN, PASSWORD, database = DB)

    client.sqlExec("""
        CREATE TABLE IF NOT EXISTS example (
            uniqueID VARCHAR[64], 
            value VARCHAR[32],
            created TIMESTAMP,
            PRIMARY KEY(uniqueID)
        );""")
        
    client.sqlExec("""
        CREATE TABLE IF NOT EXISTS related (
            id INTEGER AUTO_INCREMENT, 
            uniqueID VARCHAR[64], 
            relatedValue VARCHAR[32],
            PRIMARY KEY(id)
        );""")

    uid1 = str(uuid4())
    uid2 = str(uuid4())
    params = {
        "uid1": uid1,
        "uid2": uid2
    }
    
    resp = client.sqlExec("""
        BEGIN TRANSACTION;

        INSERT INTO example (uniqueID, value, created) 
            VALUES (@uid1, 'test1', NOW()), (@uid2, 'test2', NOW());
        INSERT INTO related (uniqueID, relatedValue) 
            VALUES (@uid1, 'related1'), (@uid2, 'related2');
        INSERT INTO related (uniqueID, relatedValue) 
            VALUES (@uid1, 'related3'), (@uid2, 'related4');

        COMMIT;
    """, params)
    
    transactionId = resp.txs[0].header.id

    result = client.sqlQuery("""
        SELECT 
            related.id,
            related.uniqueID, 
            example.value, 
            related.relatedValue, 
            example.created
        FROM related 
        JOIN example 
            ON example.uniqueID = related.uniqueID;
    """)
    for item in result:
        id, uid, value, relatedValue, created = item
        print("ITEM", id, uid, value, relatedValue, created.isoformat())

    
    result = client.sqlQuery(f"""
        SELECT 
            related.id,
            related.uniqueID, 
            example.value, 
            related.relatedValue, 
            example.created
        FROM related BEFORE TX {transactionId} 
        JOIN example BEFORE TX {transactionId} 
            ON example.uniqueID = related.uniqueID;
    """)
    print(result) # You can't see just added entries,
                  # my fellow time traveller


    # interactive session
    with client.openManagedSession(LOGIN, PASSWORD, database = DB) as session:
        transaction = session.newTx()
        for _ in range(3):
            uidNow1 = str(uuid4())
            uidNow2 = str(uuid4())
            transaction.sqlExec("""INSERT INTO example (uniqueID, value, created) 
                VALUES (@uid1, 'test130', NOW()), (@uid2, 'test131', NOW());""", {"uid1": uidNow1, "uid2": uidNow2})
        resp = transaction.commit()

        transaction = session.newTx()
        resultsSize = getResultsList(transaction)
        print("RESULTS SIZE", resultsSize) # Results size is +6 becasue we commited new transaction

        for _ in range(3):
            uidNow1 = str(uuid4())
            uidNow2 = str(uuid4())
            transaction.sqlExec("""INSERT INTO example (uniqueID, value, created) 
                VALUES (@uid1, 'test130', NOW()), (@uid2, 'test131', NOW());""", {"uid1": uidNow1, "uid2": uidNow2})

        resultsSize = getResultsList(transaction)
        print("RESULTS SIZE", resultsSize)  # Results size is +6 becasue we still not commited transaction

        transaction.rollback()

        transaction = session.newTx()

        resultsSize = getResultsList(transaction)  
        print("RESULTS SIZE", resultsSize) # Results size is -6 becasue we rollback transaction
        
def getResultsList(transaction):
    results = transaction.sqlQuery(f"""
        SELECT example.value FROM example""")
    return len(results)


if __name__ == "__main__":
    main()