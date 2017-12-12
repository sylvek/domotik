#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
from pymongo import MongoClient
import time

parser = argparse.ArgumentParser(description='store sensors, measures and triggers data to mongodb')
parser.add_argument('mqtt_hostname', metavar='mqtt_hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('mqtt_port', metavar='mqtt_port', help='port of mqtt server', nargs='?', default="1883")
parser.add_argument('mongodb_hostname', metavar='mongodb_hostname', help='hostname of mongodb server', nargs='?', default="0.0.0.0")
parser.add_argument('mongodb_port', metavar='mongodb_port', help='port of mongodb server', nargs='?', default="27017")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client_mqtt.subscribe("sensors/#")
    client_mqtt.subscribe("measures/#")

def on_message(client, userdata, msg):
    pattern = msg.topic.split("/");
    timestamp = int(time.time())
    collection = pattern[0]
    sensor = pattern[1]
    type = pattern[2]

    try:
    	value = float(msg.payload)
    except (ValueError, TypeError):
    	value = msg.payload

    message = {"sensor":sensor,"type":type,"value":value,"timestamp":timestamp}
    client_mongodb.domotik[collection].insert_one(message)

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client_mongodb.close()
    client_mqtt.disconnect()

signal.signal(signal.SIGINT, signal_handler)

client_mongodb = MongoClient(args.mongodb_hostname, int(args.mongodb_port))

client_mqtt = mqtt.Client()
client_mqtt.on_connect = on_connect
client_mqtt.on_message = on_message

client_mqtt.connect(args.mqtt_hostname, int(args.mqtt_port), 60)
client_mqtt.loop_forever()
