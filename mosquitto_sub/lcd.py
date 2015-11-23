#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import serial
import sys
import errno

parser = argparse.ArgumentParser(description='interact with lcd by sending text over serial usb')
parser.add_argument('usbport', metavar='usbport', help='usb port like /dev/ttyUSB0', nargs='?', default='/dev/ttyUSB0')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def on_connect(client, userdata, flags, rc):
    client.subscribe("triggers/lcd/text")

def on_message(client, userdata, msg):
    ser.write(msg.payload)

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    ser.close()
    client.disconnect()

try:
    signal.signal(signal.SIGINT, signal_handler)
    ser = serial.Serial(port=args.usbport, baudrate=9600, timeout=3)
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_message = on_message
    client.connect(args.hostname, int(args.port), 60)
    client.loop_forever()
except:
	sys.exit(errno.EIO)
    
