#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import struct
import RPi.GPIO as GPIO

parser = argparse.ArgumentParser(description='interact with led RGB by sending hex. rgb code')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client.subscribe("triggers/led/#")

def on_message(client, userdata, msg):
    rgb = struct.unpack('BBB', msg.payload.decode('hex'))
    rgb = [(x / 255.0) * 100 for x in rgb] # Convert 0-255 range to 0-100.
    RED.ChangeDutyCycle(rgb[0])
    GREEN.ChangeDutyCycle(rgb[1])
    BLUE.ChangeDutyCycle(rgb[2])

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    GPIO.cleanup()
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

GPIO.setmode(GPIO.BCM)
GPIO.setwarnings(False)

GPIO.setup(7, GPIO.OUT)
RED = GPIO.PWM(7, 100)
RED.start(0)

GPIO.setup(8, GPIO.OUT)
GREEN = GPIO.PWM(8, 100)
GREEN.start(0)

GPIO.setup(25, GPIO.OUT)
BLUE = GPIO.PWM(25, 100)
BLUE.start(0)

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
