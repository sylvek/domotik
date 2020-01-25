import paho.mqtt.client as mqtt
import argparse
import time
import sys
import errno
import signal

parser = argparse.ArgumentParser(description='send celsus temperature every minute on topic sensors/esp9266/temp')
parser.add_argument('temp_sensor', metavar='temp_sensor', help='path to 1-wire temperature sensor', nargs='?', default='/sys/bus/w1/devices/28-0315820a56ff/w1_slave')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

temp_sensor = args.temp_sensor
run = True

def signal_handler(signal, frame):
    global run
    print "Ending and cleaning up"
    run = False

try:
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    client = mqtt.Client()
    client.connect(args.hostname, int(args.port), 120)
    client.loop_start()
except Exception as e:
    print e
    sys.exit(errno.EIO)

def temp_raw():
    f=open(temp_sensor, 'r')
    lines = f.readlines()
    f.close()
    return lines

def read_temp():
    lines = temp_raw()
    while lines[0].strip()[-3:] != 'YES':
        time.sleep(0.2)
        lines = temp_raw()
    temp_output = lines[1].find('t=')
    if temp_output != -1:
        temp_string = lines[1].strip()[temp_output+2:]
        temp_c = float(temp_string) / 1000.0
        return round(temp_c, 2)

while run:
    try:
        temp = read_temp()
        if temp < 50:
            client.publish("sensors/esp8266/temp", temp)
        time.sleep(60)
    except:
        pass

client.disconnect()
