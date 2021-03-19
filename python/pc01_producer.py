#!/usr/bin/env python
from immudb.client import ImmudbClient
import json
ic=ImmudbClient()
ic.login(username="immudb", password="immudb")


key="BANK_TRANSFER_001"
data={
    "bank_transaction_id":762349991002,
    "from":"Bob",
    "to":"Mark",
    "amount":420.00
    }

print("Entering money transfer data into immudb")

# values must encode to bytes
ret=ic.verifiedSet( key.encode("utf8"), json.dumps(data).encode("utf8") )
print("Data saved with transaction {}, verified: {}".format(ret.id, ret.verified))
