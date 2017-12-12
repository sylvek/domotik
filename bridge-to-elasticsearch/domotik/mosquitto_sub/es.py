#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
from elasticsearch import Elasticsearch
from datetime import datetime

parser = argparse.ArgumentParser(description='store sensors, measures and triggers data to mongodb')
parser.add_argument('mqtt_hostname', metavar='mqtt_hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('mqtt_port', metavar='mqtt_port', help='port of mqtt server', nargs='?', default="1883")
parser.add_argument('elasticsearch_hostname', metavar='elasticsearch_hostname', help='hostname of elasticsearch server', nargs='?', default="0.0.0.0")
parser.add_argument('elasticsearch_port', metavar='elasticsearch_port', help='port of elasticsearch server', nargs='?', default="27017")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client_mqtt.subscribe("sensors/#")
    client_mqtt.subscribe("measures/#")

def on_message(client, userdata, msg):
    pattern = msg.topic.split("/");
    timestamp = datetime.now()
    collection = pattern[0]
    sensor = pattern[1]
    type = pattern[2]

    try:
    	value = float(msg.payload)
    except (ValueError, TypeError):
    	value = msg.payload

    message = {"sensor":sensor,"type":type,"value":value,"timestamp":timestamp}
    client_elasticsearch.index(index="domotik", doc_type=collection, body=message)

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client_elasticsearch.close()
    client_mqtt.disconnect()

signal.signal(signal.SIGINT, signal_handler)

client_elasticsearch = Elasticsearch([args.elasticsearch_hostname], port=args.elasticsearch_port)

client_mqtt = mqtt.Client()
client_mqtt.on_connect = on_connect
client_mqtt.on_message = on_message

client_mqtt.connect(args.mqtt_hostname, int(args.mqtt_port), 60)
client_mqtt.loop_forever()
