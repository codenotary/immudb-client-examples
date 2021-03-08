#!/usr/bin/env python

from immudb.client import ImmudbClient

ic = ImmudbClient()
ic.login(username="immudb", password="immudb")

key = "Hello".encode('utf8')
value = "Immutable World!".encode('utf8')

# set a key/value pair
ic.set(key, value)

# reads back the value
readback = ic.get(key)
saved_value = readback.value.decode('utf8')
print("Hello", saved_value)
