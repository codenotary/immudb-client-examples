from immudb.client import ImmudbClient
import string,random
import itertools
import time

SIZE=5000000
CHUNKSIZE=5000

def chunked(it, size):
    it = iter(it)
    while True:
        p = dict(itertools.islice(it, size))
        if not p:
            break
        yield p

ic = ImmudbClient("localhost:3322")
ic.login("immudb","immudb")

print("Preparing dictionary:")
# let's fill a big dictionary:
big_dict={}
for i in range(0,SIZE):
    big_dict["verymassif:{:08X}".format(i).encode('utf8')]="value:{:08f}".format(random.random()).encode('utf8')
    if (i%CHUNKSIZE)==0:
        print("\r{:02.1f}%".format(i*100.0/SIZE),end='')
    
print("\nDone\nInserting {} values:".format(SIZE))

# now we put all the key/value pairs in immudb
written=0
t0=time.time()
for chunk in chunked(big_dict.items(), CHUNKSIZE):
    response=ic.setAll(chunk)
    # the response holds the new index position of the merkele tree
    assert type(response)!=int
    written+=CHUNKSIZE
    print("\r{:02.1f}%".format(written*100.0/SIZE),end='')
t1=time.time()

print("\nDone")
print("{} keys written in {:3.2f} seconds".format(SIZE,t1-t0))

