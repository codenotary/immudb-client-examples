#/bin/bash

echo -n immudb | immuclient login immudb

# Set a key "hello" to the value "immutable world"
immuclient set "hello" "immutable world"

# get the value for the key "hello"
immuclient get "hello"

# set and verify the key "welcome" with value "immudb"
immuclient safeset "welcome" "immudb"

# Retrieve and verify the entry with key "welcome"
immuclient safeget "welcome"

immuclient logout
