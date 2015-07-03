#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse

parser = argparse.ArgumentParser(description='calculate mean of power consumption')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

count = 0
sum = 0

def on_connect(client, userdata, flags, rc):
    client.subscribe("sensors/cc128/watt")

def on_message(client, userdata, msg):
    global count
    global sum
    count += 1
    sum += int(msg.payload)
    client.publish("sensors/cc128/watt/mean", (sum/count))

client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
