#!/usr/bin/python
import pymongo
import argparse
import signal
import datetime

parser = argparse.ArgumentParser(description='(re)generate measures for meanPerHour')
parser.add_argument('hostname', metavar='hostname', help='hostname of mongodb server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mongodb server', nargs='?', default="27017")
args = parser.parse_args()


def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.close()

signal.signal(signal.SIGINT, signal_handler)

currentHour = -1
sum = 0
count = 0

client = pymongo.MongoClient(args.hostname, int(args.port))
db = client.domotik
for doc in db.sensors.find({ 'type':'watt' }).sort([('timestamp', pymongo.ASCENDING)]):
    timestamp = doc['timestamp']
    hour = datetime.datetime.fromtimestamp(timestamp).hour
    sum += doc['value']
    count += 1

    if (currentHour is -1):
        currentHour = hour

    if (hour is not currentHour):
        print "{timestamp:" + str(timestamp) + ", sensor:'meanPerHour', type:'watt', value:" + str(sum/count) + "}"
        sum = 0
        count = 0
        currentHour = hour
