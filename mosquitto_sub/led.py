#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import struct
import RPi.GPIO as GPIO
import time

parser = argparse.ArgumentParser(description='interact with led RGB by sending hex. rgb code')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client.subscribe("triggers/led/+")

def on_message(client, userdata, msg):
    topic = msg.topic
    rgb = msg.payload
    if topic == "triggers/led/update":
        update_led_value(rgb)
    if topic == "triggers/led/blink":
        blink_led_value(rgb)

def blink_led_value(value):
    for x in xrange(0, 3):
        update_led_value(value)
        time.sleep(1.000)
        update_led_value("ffffff")
        time.sleep(1.000)

def update_led_value(value):
    rgb = struct.unpack('BBB', value.decode('hex'))
    rgb = [(x / 255.0) * 100 for x in rgb] # Convert 0-255 range to 0-100.
    change_duty_cycle(rgb)

def change_duty_cycle(rgb):
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
RED.start(100)

GPIO.setup(8, GPIO.OUT)
GREEN = GPIO.PWM(8, 100)
GREEN.start(100)

GPIO.setup(25, GPIO.OUT)
BLUE = GPIO.PWM(25, 100)
BLUE.start(100)

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
