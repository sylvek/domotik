#!/usr/bin/env python
import socket
import time
import sys
import argparse

parser = argparse.ArgumentParser(description='send data via pipeline to graphite')
parser.add_argument('hostname', metavar='hostname', help='hostname of graphite server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of graphite server', nargs='?', default="2003")
args = parser.parse_args()

for line in sys.stdin:
		timestamp = int(time.time())
		message = '%s %d\n' % (line.rstrip(), timestamp)

		sock = socket.socket()
		sock.connect((args.hostname, int(args.port)))
		sock.sendall(message)
		sock.close()
