#!/usr/bin/env python3
from immudb import ImmudbClient
import string,random, argparse

def get_random_string(length):
	return ''.join(random.choice(string.ascii_letters+string.digits) for i in range(length))

def NotarizeTransaction(ic, transaction_id, sender, receiver, amount, description):
	key='TRANSACTION:{}:SENDER'.format(transaction_id)
	ret=ic.verifiedSet(key.encode('utf8'), sender.encode('utf8'))
	assert ret.verified

	key='TRANSACTION:{}:RECEIVER'.format(transaction_id) 
	ret=ic.verifiedSet(key.encode('utf8'), receiver.encode('utf8'))
	assert ret.verified

	key='TRANSACTION:{}:AMOUNT'.format(transaction_id) 
	ret=ic.verifiedSet(key.encode('utf8'), str(amount).encode('utf8'))
	assert ret.verified

	key='TRANSACTION:{}:DESCRIPTION'.format(transaction_id) 
	ret=ic.verifiedSet(key.encode('utf8'), description.encode('utf8'))
	assert ret.verified

def NotarizeReceipit(ic, transaction_id, receipit_filename):
	with open(receipit_filename,"rb") as f:
		receipit_content=f.read()
	key='TRANSACTION:{}:RECEIPIT'.format(transaction_id) 
	ret=ic.verifiedSet(key.encode('utf8'), receipit_content)
	assert ret.verified
	

def NotarizeTransactionBatch(ic, transaction_id, sender, receiver, amount, description):
	transaction={
		'TRANSACTION:{}:SENDER'.format(transaction_id) : sender,
		'TRANSACTION:{}:RECEIVER'.format(transaction_id) : receiver,
		'TRANSACTION:{}:AMOUNT'.format(transaction_id) : str(amount),
		'TRANSACTION:{}:DESCRIPTION'.format(transaction_id) : description,
		}
	encoded_transaction={
		x[0].encode('utf8'):x[1].encode('utf8') for x in transaction.items()
		}
	ret=ic.setAll(encoded_transaction)
	ret=ic.verifiedTxById(ret.id)
	for i in transaction.keys():
		assert i.encode('utf8') in ret

parser = argparse.ArgumentParser(description='Showcase demo')
parser.add_argument('--sender', type=str, default='Tim', help='Sender Name')
parser.add_argument('--recipient', type=str, default='Tom', help='Recipient Name')
parser.add_argument('--amount', type=float, default=500.0, help='Amount transferred')
parser.add_argument('--description', type=str, default='payment', help='Transaction description')
parser.add_argument('--receipit', type=str, help='Recipient filename')
parser.add_argument('--batched', default=False, action='store_true', help='Use batch operation')

args=parser.parse_args()

ic=ImmudbClient()
ic.login("immudb","immudb")
ic.databaseUse("defaultdb")

transaction_id=get_random_string(16)
print("Transaction id:",transaction_id)
if args.batched:
	NotarizeTransactionBatch(ic, transaction_id, args.sender, args.recipient, args.amount, args.description)
else:
	NotarizeTransaction(ic, transaction_id, args.sender, args.recipient, args.amount, args.description)
if args.receipit!=None:
	NotarizeReceipit(ic, transaction_id, args.receipit)
