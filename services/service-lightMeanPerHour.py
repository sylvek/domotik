#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import struct

parser = argparse.ArgumentParser(description='blink the led with a certain color one time per hour')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('trigger_out', metavar='trigger_out', help='trigger one time per hour')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

previous_value = 0

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    global previous_value
    current_value = int(msg.payload)
    if (previous_value > 0):
        percent = limit_percent(current_value / previous_value)
        red = int(255 - 125 * percent)
        green = int(125 * percent)
        color = struct.pack('BBB',*(red,green,255)).encode('hex')
        client.publish(args.trigger_out, color)
    previous_value = current_value

def limit_percent(percent):
    if (percent > 2.0):
        return 2.0

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
