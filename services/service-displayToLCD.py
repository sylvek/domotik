#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import struct

parser = argparse.ArgumentParser(description='forwards message to trigger after formatting it')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('trigger_out', metavar='trigger_out', help='trigger message')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    topic = msg.topic.split('/')
    payload = msg.payload
    unit = "degree" if topic[2] == "temp" else topic[2]
    sensor = "salon" if topic[1] == "cc128" else "dehors c." if topic[1] == "esp12e" else "dehors b." if topic[1] == "thn132n" else "salle de bain" if topic[1] == "esp8266" else topic[1]
    client.publish(args.trigger_out, sensor + ";" + payload + ";" + unit)


def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
