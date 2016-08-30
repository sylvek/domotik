#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import datetime
import json
import os.path

parser = argparse.ArgumentParser(description='try to discover the hot water consumption by analysing when the hot water tank goes off. That should occur only one time per night.')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('trigger_out', metavar='trigger_out', help='triggered topic')
parser.add_argument('percent', metavar='percent', help='trigger limit in percent between two values. 0.20 means that the new measure looses 80p of value.', nargs='?', default="0.20")
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

previous_value = 0
trigger = False

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    global previous_value
    global trigger
    global day
    now = datetime.datetime.now()
    current_value = float(msg.payload)
    current_time_in_minute = now.hour * 60 + now.minute # minutes since the begining of current day
    p = float(args.percent)
    if (not trigger and current_time_in_minute > 1350): # 1350 minutes means 22h30
        trigger = True
        previous_value = current_value
    if (trigger and current_value < previous_value * p):
        trigger = False
        elapsed_time_in_minute = (1350 - current_time_in_minute) if now.day == day else (90 + current_time_in_minute) # 90 means minutes between 22h30 and midnight
        client.publish(args.trigger_out, elapsed_time_in_minute)

def signal_handler(signal, frame):
    with open(__file__ + "." + args.service_name + ".previous", 'w') as outfile:
        global previous_value
        global day
        global trigger
        json.dump({'day': day, 'previous_value': previous_value, 'trigger': trigger}, outfile)
    print "Ending and cleaning up"
    client.disconnect()

day = datetime.datetime.now().day

if os.path.exists(__file__ + "." + args.service_name + ".previous"):
    with open(__file__ + "." + args.service_name + ".previous", 'r') as infile:
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
