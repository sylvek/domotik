#!/usr/bin/python
from BaseHTTPServer import BaseHTTPRequestHandler,HTTPServer
import paho.mqtt.client as mqtt
import signal
import argparse

parser = argparse.ArgumentParser(description='trigg an event when web server receive a request')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('trigger_out', metavar='trigger_out', help='trigg out when something is received')
parser.add_argument('http_port', metavar='http_port', help='web server port', nargs='?', default="8888")
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

run = True

class myHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        client.publish(args.trigger_out, self.client_address[0])
        self.send_response(200)
        self.send_header('Content-type','text/html')
        self.end_headers()
        return

def signal_handler(signal, frame):
    global run
    print "Ending and cleaning up"
    run = False
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.connect(args.hostname, int(args.port), 60)

server = HTTPServer(('', int(args.http_port)), myHandler)
server.timeout = 1

while run:
    client.loop()
    server.handle_request()

server.server_close()
