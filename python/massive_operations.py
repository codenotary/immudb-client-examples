
from immudb.client import ImmudbClient
import string
import random


def get_random_string(length):
    return ''.join(random.choice(string.printable) for i in range(length))


ic = ImmudbClient()
ic.login("immudb", "immudb")

# let's fill a big dictionary:
xset = {}
for i in range(0, 1000):
    xset["massif:{:04X}".format(i).encode(
        'utf8')] = get_random_string(32).encode('utf8')

# now we put all the key/value pairs in immudb
response = ic.setAll(xset)

# the response holds the new index position of the merkele tree
assert type(ic.setAll(xset)) != int

# let's read back all the values in another dictionary,
# and check the values
yset = ic.getAll(xset.keys())
for i in yset.keys():
    assert xset[i] == yset[i]
