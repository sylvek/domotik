#!/usr/bin/env python
import argparse
import sys
import urllib
import time

parser = argparse.ArgumentParser(description='send data via pipeline to thinkspeak')
parser.add_argument('hostname', metavar='key', help='channel api key')
parser.add_argument('port', metavar='field', help='field id')
args = parser.parse_args()

while sys.stdin:
	line=sys.stdin.readline().rstrip()
	if line:
            url = "https://api.thingspeak.com/update?api_key=" + args.hostname + "&field" + args.port + "=" + line
            response = urllib.urlopen(url);
	else:
            time.sleep(1)
