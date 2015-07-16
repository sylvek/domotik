#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import serial
import sys
import errno
import re

parser = argparse.ArgumentParser(description='fetch data from cc128 and push it to mqtt')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

system = 'envi' #'classic'
usbport = '/dev/ttyUSB0'

def sensor_data_check( sensor, watt, tmpr ):
    sensor = int(sensor)
    watt = int(watt)
    tmpr = float(tmpr)
    client.publish("sensors/cc128/watt", watt)
    client.publish("sensors/cc128/temp", tmpr)

try:
	ser = serial.Serial(port=usbport, baudrate=57600, bytesize=serial.EIGHTBITS, parity=serial.PARITY_NONE, stopbits=serial.STOPBITS_ONE, timeout=3)
	client = mqtt.Client()
	client.connect(args.hostname, int(args.port), 60)
except:
	sys.exit(errno.EIO)

while True:
	try:
		line = ser.readline()
		line = line.rstrip('\r\n')
		clamps = False

		if system == 'classic':
			r = re.search(r"<ch1><watts>(\d+)<\/watts><\/ch1><ch2><watts>(\d+)<\/watts><\/ch2><ch3><watts>(\d+)<\/watts><\/ch3><tmpr>(.+?)</tmpr>", line)
			if r:
				tmpr = r.group(4)
				s = 1
				clamp1 = r.group(1)
				clamp2 = r.group(2) if int(r.group(2)) > 0 else False
				clamp3 = r.group(3) if int(r.group(3)) > 0 else False
		else:
			r = re.search(r"<tmpr>(.+?)</tmpr><sensor>(\d)+</sensor>.+<ch1><watts>(\d+)<\/watts><\/ch1>(<ch2><watts>(\d+)<\/watts><\/ch2>)?(<ch3><watts>(\d+)<\/watts><\/ch3>)?", line)
			if r:
				tmpr = r.group(1)
				s = r.group(2)
				clamp1 = r.group(3)
				clamp2 = r.group(5) if r.group(5) else False
				clamp3 = r.group(7) if r.group(7) else False

		if r:
			watt_sum = int(clamp1)
			# more than 1 clamp
			if clamp2:
				sensor = int('2'+str(s))
				if sensors and sensors.has_key(sensor):
					watt = int(clamp2)
					watt_sum += watt
					clamps = True

			if clamp3:
				sensor = int('3'+str(s))
				if sensors and sensors.has_key(sensor):
					watt = int(clamp3)
					watt_sum += watt
					clamps = True

			if clamps:
				sensor = int('1'+str(s))
				watt = int(clamp1)

			sensor_data_check( s, watt_sum, tmpr )
	except:
		sys.exit(errno.EIO)

client.loop_forever()
