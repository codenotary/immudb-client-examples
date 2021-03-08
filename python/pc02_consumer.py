#!/usr/bin/env python
from immudb.client import ImmudbClient
import json
ic=ImmudbClient()
ic.login(username="immudb", password="immudb")


key="BANK_TRANSFER_001"

print("Reading money transfer data from immudb")

# values must encode to bytes
ret=ic.verifiedGet( key.encode("utf8"))
print("Data was saved with transaction {}, verified: {}".format(ret.id, ret.verified))
data=json.loads(ret.value.decode('utf8'))
print("Data :",data)
