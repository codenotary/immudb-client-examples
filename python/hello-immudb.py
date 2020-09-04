#!/usr/bin/env python
from immu.client import ImmuClient

client = ImmuClient("localhost:3322")
client.login("immudb", "immudb")
client.safeSet(b"hello", b"world")
client.safeGet(b"hello")
client.logout()
