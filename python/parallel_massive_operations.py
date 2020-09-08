from immudb.client import ImmudbClient
import string,random
import itertools
import time
import multiprocessing 

SIZE=1000000
CHUNKSIZE=1000

def chunked(it, size):
    it = iter(it)
    while True:
        p = dict(itertools.islice(it, size))
        if not p:
            break
        yield p


def massive_test(taskid:int):
    ic = ImmudbClient("localhost:3322")
    ic.login("immudb","immudb")

    # let's fill a big dictionary:
    big_dict={}
    for i in range(0,SIZE):
        big_dict["verymassif:{:08X}".format(i).encode('utf8')]="value:{:08f}".format(random.random()).encode('utf8')
  

    # now we put all the key/value pairs in immudb
    written=0
    t0=time.time()
    for chunk in chunked(big_dict.items(), CHUNKSIZE):
        response=ic.setAll(chunk)
        # the response holds the new index position of the merkele tree
        assert type(response)!=int
        written+=CHUNKSIZE
    t1=time.time()
    print("TASK{}:  {} keys written in {:3.2f} seconds".format(taskid,SIZE,t1-t0))
    return t1-t0

plist=[]
for i in range(0,4):
    p=multiprocessing.Process(target=massive_test, args=(i,))
    p.start()
    plist.append(p)
for p in plist:
    p.join()
