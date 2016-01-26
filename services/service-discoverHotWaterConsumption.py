#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import datetime
import json
import os.path

parser = argparse.ArgumentParser(description='try to discover the hot water consumption by analysing when the hot water tank goes off. That should occur only one time per night.')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('trigger_out', metavar='trigger_out', help='trigger message')
parser.add_argument('percent', metavar='percent', help='trigger limit in percent between two values. 0.10 => 10percent means that the new measure looses 90p of value.', nargs='?', default="0.10")
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

previous_value = 0
trigger = False

minutes_previous_midnight = 90

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    global previous_value
    global trigger
    global day
    now = datetime.datetime.now()
    current_value = float(msg.payload)
    current_day = now.day
    p = float(args.percent)
    if (current_day is not day):
        trigger = True
        day = current_day
    if (trigger and current_value < previous_value * p):
        current_time_in_minute = now.hour * 60 + now.minute
        client.publish(args.trigger_out, current_time_in_minute + minutes_previous_midnight)
        trigger = False
    previous_value = current_value

def signal_handler(signal, frame):
    with open(__file__ + ".previous", 'w') as outfile:
        global previous_value
        global day
        global trigger
        json.dump({'day': day, 'previous_value': previous_value, 'trigger': trigger}, outfile)
    print "Ending and cleaning up"
    client.disconnect()

day = datetime.datetime.now().day

if os.path.exists(__file__ + ".previous"):
    with open(__file__ + ".previous", 'r') as infile:
        previous = json.load(infile)
        day = previous['day']
        previous_value = previous['previous_value']
        trigger = previous['trigger']

signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
