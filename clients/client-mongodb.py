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
	line=sys.stdin.readline().rstrip().split();
	if line:
		pattern = line[0].split("/");
		timestamp = int(time.time())
		collection = pattern[0]
		sensor = pattern[1]
		value = str(line[1]) if isinstance(line[1], str) else float(line[1])

		message = {"sensor":sensor,"type":type,"value":value,"timestamp":timestamp}
		db = client.domotik
		db[collection].insert_one(message)
	else:
	    	time.sleep(1)
