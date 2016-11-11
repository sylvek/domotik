#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import os
import time

parser = argparse.ArgumentParser(description='checks sensor availability')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('delay', metavar='delay', help='delay without two measures to alert (in min)', nargs='?', default="30")
parser.add_argument('trigger_out', metavar='trigger_out', help='triggered if elapsed time is greater than delay')
parser.add_argument('trigger_value', metavar='trigger_value', help='triggered value')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

previous_value = time.time()
run = True

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    global previous_value
    previous_value = time.time()

def signal_handler(sig, frame):
    if sig is not signal.SIGUSR1:
        print "Ending and cleaning up"
        global run
        run = False
        client.disconnect()

def alert_if_needed():
    global previous_value
    now = time.time()
    if (previous_value + int(args.delay) * 60 < now):
        client.publish(args.trigger_out, args.trigger_value)
        previous_value = now

signal.signal(signal.SIGUSR1, signal_handler)
signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)

while run:
    client.loop()
    alert_if_needed()
