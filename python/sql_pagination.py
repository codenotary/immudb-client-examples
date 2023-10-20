from immudb.client import ImmudbClient
import uuid
MAX_N = 2000
def main():
    client = ImmudbClient("localhost:3322")
    client.login("immudb", "immudb", "defaultdb")
    client.sqlExec("""CREATE TABLE IF NOT EXISTS pagination( ID INTEGER AUTO_INCREMENT, VALUE VARCHAR[255], PRIMARY KEY(ID) )""")

    for index in range(0, MAX_N):
        print("Adding records...", index, f"{round((index / MAX_N) * 100, 2)}%")
        value = str(uuid.uuid4())
        client.sqlExec("INSERT INTO pagination(VALUE) VALUES(@someValue)", {"someValue": value})

    query = """SELECT ID, VALUE FROM pagination ORDER BY ID LIMIT 999"""

    # Pagination method
    wholeResult = []
    result = client.sqlQuery(query)
    while result:
        wholeResult.extend(result)
        queryNext = """SELECT ID, VALUE FROM pagination WHERE id > @lastId ORDER BY ID LIMIT 999"""
        lastId = result[-1][0]
        result = client.sqlQuery(queryNext, {"lastId": lastId})
    print(len(wholeResult))

    # Offset method - less performance
    wholeResult = []
    result = client.sqlQuery(query)
    offset = 0
    query = f"""SELECT ID, VALUE FROM pagination ORDER BY ID LIMIT 999"""
    while result:
        offset = offset + len(result)
        wholeResult.extend(result)
        queryNext = f"""SELECT ID, VALUE FROM pagination ORDER BY ID LIMIT 999 OFFSET {offset}"""
        lastId = result[-1][0]
        result = client.sqlQuery(queryNext)
    print(len(wholeResult))

main()