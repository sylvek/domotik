#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import syslog

parser = argparse.ArgumentParser(description='log sensors, measures and triggers data to syslog')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client.subscribe("sensors/#")
    client.subscribe("triggers/#")
    client.subscribe("measures/#")

def on_message(client, userdata, msg):
    syslog.syslog(msg.topic + " " + str(msg.payload)[:80])

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
