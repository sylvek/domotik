#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import serial
import sys
import errno
import re
import signal

parser = argparse.ArgumentParser(description='fetch data from thn132n and push it to mqtt')
parser.add_argument('usbport', metavar='usbport', help='usb port like /dev/ttyUSB0', nargs='?', default='/dev/ttyUSB0')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

usbport = args.usbport
run = True

def sensor_data_check( tmpr, moisture, battery ):
    moisture = int(moisture)
    battery = int(battery)
    tmpr = float(tmpr)
    client.publish("sensors/thn132n/battery", battery)
    client.publish("sensors/thn132n/moisture", moisture)
    client.publish("sensors/thn132n/temp", tmpr)

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    ser.close()
    client.disconnect()
    run = False

try:
	signal.signal(signal.SIGINT, signal_handler)
	ser = serial.Serial(port=usbport, baudrate=9600, timeout=3)
	client = mqtt.Client()
	client.connect(args.hostname, int(args.port), 60)
except:
	sys.exit(errno.EIO)

while run:
	try:
		client.loop()
		line = ser.readline()
        	line = line.rstrip('\r\n')
        	if ( len(line) > 0 ):
            		line = line.split(',',3)
            		sensor_data_check( line[0], line[1], line[2] )
	except:
		sys.exit(errno.EIO)
