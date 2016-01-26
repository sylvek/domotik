#!/usr/bin/python
from __future__ import division
import paho.mqtt.client as mqtt
import argparse
import signal
import datetime

parser = argparse.ArgumentParser(description='calculate a mean and send a measure one time per hour')
parser.add_argument('sensor_in', metavar='sensor_in', help='sensor path given')
parser.add_argument('sensor_out', metavar='sensor_out', help='sensor path given a resulted mean')
parser.add_argument('measure_out', metavar='measure_out', help='measure path given a resulted mean (one time per hour)')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

count = 0
sum = 0

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.sensor_in)

def on_message(client, userdata, msg):
    global count
    global sum
    global hour
    count += 1
    sum += int(msg.payload)
    mean = sum/count
    client.publish(args.sensor_out, mean)
    currentHour = datetime.datetime.now().hour
    if (currentHour is not hour):
        count = 1
        sum = mean
        hour = currentHour
        client.publish(args.measure_out, round(mean, 2))

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.disconnect()

hour = datetime.datetime.now().hour

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
