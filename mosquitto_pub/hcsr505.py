#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import struct
import RPi.GPIO as GPIO
import time
import datetime

parser = argparse.ArgumentParser(description='interact with the hcsr505 by reading the motion value')
parser.add_argument('delay', metavar='delay', help='delay in seconds between 2 event', nargs='?', default="15")
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

run = True

def signal_handler(signal, frame):
    global run
    print "Ending and cleaning up"
    GPIO.remove_event_detect(14)
    GPIO.cleanup()
    client.disconnect()
    run = False

def event_detected(channel):
    global locked_datetime
    value = GPIO.input(14)
    current_datetime = datetime.datetime.now()
    trigger_datetime = locked_datetime + datetime.timedelta(seconds=int(args.delay))
    if trigger_datetime < current_datetime:
        client.publish("sensors/hcsr505/event", value)
        locked_datetime = current_datetime

try:
    signal.signal(signal.SIGINT, signal_handler)
    client = mqtt.Client()

    locked_datetime = datetime.datetime.now()

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
