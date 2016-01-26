#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import struct

parser = argparse.ArgumentParser(description='do something if measure of sensor is less than the limit')
parser.add_argument('sensor_in', metavar='sensor_in', help='sensor path given')
parser.add_argument('trigger_out', metavar='trigger_out', help='trigger out, ex: blink a led')
parser.add_argument('limit', metavar='limit', help='limit')
parser.add_argument('trigger_value', metavar='trigger_value', help='triggered value')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.sensor_in)

def on_message(client, userdata, msg):
    current_value = float(msg.payload)
    limit = int(args.limit)
    if (current_value < limit):
        client.publish(args.trigger_out, args.trigger_value)

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
