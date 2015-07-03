#!/usr/bin/python
import urllib, json
import argparse

parser = argparse.ArgumentParser(description='retrieve data from measureit (cc128 project)')
parser.add_argument('hostname', metavar='hostname', help='hostname of graphite server')
parser.add_argument('sensor', metavar='sensor', help='tmpr or watt')
args = parser.parse_args()

url = args.hostname + "/php/measureit_functions.php?do=summary_start"
response = urllib.urlopen(url);
data = json.loads(response.read())
print data[0][args.sensor] #tmpr or watt
