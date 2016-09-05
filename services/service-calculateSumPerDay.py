#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import datetime
import json
import os.path

parser = argparse.ArgumentParser(description='calculate a sum of measures per day')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('measure_out', metavar='measure_out', help='measure path given a resulted sum (one time per day)')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

sum = 0

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    global sum
    global day
    sum += float(msg.payload)
    currentDay = datetime.datetime.now().day
    if (currentDay is not day):
        client.publish(args.measure_out, sum)
        sum = 0
        day = currentDay

def signal_handler(sig, frame):
    with open(__file__ + "." + args.service_name + ".previous", 'w') as outfile:
        global day
        global sum
        json.dump({'day': day, 'sum': sum}, outfile)
    if sig is not signal.SIGUSR1:
        print "Ending and cleaning up"
        client.disconnect()

day = datetime.datetime.now().day

if os.path.exists(__file__ + "." + args.service_name + ".previous"):
    with open(__file__ + "." + args.service_name + ".previous", 'r') as infile:
        previous = json.load(infile)
        day = previous['day']
        sum = previous['sum']

signal.signal(signal.SIGUSR1, signal_handler)
signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
