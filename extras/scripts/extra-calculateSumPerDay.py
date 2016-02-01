#!/usr/bin/python
import pymongo
import argparse
import signal
import datetime

parser = argparse.ArgumentParser(description='(re)generate measures for sumPerDay')
parser.add_argument('hostname', metavar='hostname', help='hostname of mongodb server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mongodb server', nargs='?', default="27017")
args = parser.parse_args()


def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.close()

signal.signal(signal.SIGINT, signal_handler)

currentDay = -1
sum = 0
count = 0

client = pymongo.MongoClient(args.hostname, int(args.port))
db = client.domotik
for doc in db.sensors.find({ 'type':'watt' }).sort([('timestamp', pymongo.ASCENDING)]):
    timestamp = doc['timestamp']
    day = datetime.datetime.fromtimestamp(timestamp).day
    sum += doc['value']
    count += 1

    if (currentDay is -1):
        currentDay = day

    if (day is not currentDay):
        print "{timestamp:" + str(timestamp) + ", sensor:'sumPerDay', type:'watt', value:" + str(sum/count*24) + "}"
        sum = 0
        count = 0
        currentDay = day
