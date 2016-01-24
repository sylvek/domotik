#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import struct
import RPi.GPIO as GPIO
import time

parser = argparse.ArgumentParser(description='interact with the hcsr505 by reading the motion value')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

run = True

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    GPIO.remove_event_detect(14)
    GPIO.cleanup()
    client.disconnect()
    run = False

def event_detected(channel):
    value = GPIO.input(14)
    client.publish("sensors/hcsr505/event", value)

try:
    signal.signal(signal.SIGINT, signal_handler)
    client = mqtt.Client()

    GPIO.setmode(GPIO.BCM)
    GPIO.setwarnings(False)

    GPIO.setup(14, GPIO.IN)
    GPIO.add_event_detect(14, GPIO.BOTH, callback=event_detected, bouncetime=100)

    client.connect(args.hostname, int(args.port), 60)
except:
	sys.exit(errno.EIO)

while run:
	try:
		client.loop()
        value = GPIO.input(14)
        client.publish("sensors/hcsr505/motion", value)
        time.sleep(10)
	except:
		sys.exit(errno.EIO)
