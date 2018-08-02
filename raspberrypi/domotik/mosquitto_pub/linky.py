#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import serial
import sys
import errno
import re
import signal
import os

parser = argparse.ArgumentParser(description='fetch data from linky and push it to mqtt')
parser.add_argument('usbport', metavar='usbport', help='usb port like /dev/serial0', nargs='?', default='/dev/serial0')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

usbport = args.usbport
run = True

def signal_handler(signal, frame):
    global run
    print "Ending and cleaning up"
    ser.close()
    client.disconnect()
    run = False

try:
    print("Starting listening Linky")
    signal.signal(signal.SIGINT, signal_handler)

    while not os.path.exists(usbport):
        print "waiting for", usbport
        time.sleep(2)

    ser = serial.Serial(port=usbport, baudrate=1200, bytesize=serial.SEVENBITS, parity=serial.PARITY_EVEN, stopbits=serial.STOPBITS_ONE, timeout=1)
    client = mqtt.Client()
    client.connect(args.hostname, int(args.port), 60)
except Exception as e:
    print e
    sys.exit(errno.EIO)

while run:
    try:
        client.loop()
        line = ser.readline()
        if line.startswith( 'PAPP' ):
                client.publish("sensors/linky/watt", line[5 : 10])
    except:
        sys.exit(errno.EIO)
