#!/usr/bin/env python
import time
import sys
import argparse
import pymongo
from pymongo import MongoClient

parser = argparse.ArgumentParser(description='send data via pipeline to mongodb')
parser.add_argument('hostname', metavar='hostname', help='hostname of mongodb server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mongodb server', nargs='?', default="27017")
args = parser.parse_args()

while sys.stdin:
	client = MongoClient(args.hostname, int(args.port))
	line=sys.stdin.readline().rstrip().replace('/', '-').split();
	if line:
		timestamp = int(time.time())
		message = {"value":float(line[1]),"timestamp":timestamp}
		db = client[line[0]]
		db.values.insert_one(message)
	else:
	    	time.sleep(1)
