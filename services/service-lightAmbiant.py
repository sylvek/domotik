#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import datetime
import schedule
import time

parser = argparse.ArgumentParser(description='light an ambiant style :)')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('trigger_out', metavar='trigger_out', help='trigger to control the light')
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

run = True

def publish(m=''):
    client.publish(args.trigger_out, m)

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    global run
    run = False
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)
client = mqtt.Client()
client.connect(args.hostname, int(args.port), 60)

# morning
schedule.every().day.at("7:00").do(publish, m='#1C3EA0')
schedule.every().day.at("7:10").do(publish, m='#2550D0')
schedule.every().day.at("7:20").do(publish, m='#2856DE')
schedule.every().day.at("7:30").do(publish, m='#2B5DF1')
schedule.every().day.at("7:40").do(publish, m='#2C60F8')
schedule.every().day.at("8:00").do(publish, m='#F8F800')
schedule.every().day.at("8:10").do(publish, m='#F87C00')
schedule.every().day.at("8:15").do(publish, m='#F800BA')
schedule.every().day.at("8:20").do(publish, m='#000000')

# night
schedule.every().day.at("19:00").do(publish, m='#98BF57')
schedule.every().day.at("19:15").do(publish, m='#33BFAE')
schedule.every().day.at("19:30").do(publish, m='#3368BF')
schedule.every().day.at("19:45").do(publish, m='#BF338B')
schedule.every().day.at("20:00").do(publish, m='#BFAE33')
schedule.every().day.at("20:30").do(publish, m='#0D3F93')
schedule.every().day.at("21:00").do(publish, m='#933593')
schedule.every().day.at("22:00").do(publish, m='#7A5921')
schedule.every().day.at("23:00").do(publish, m='#234111')
schedule.every().day.at("00:00").do(publish, m='#000000')

while run:
    client.loop()
    schedule.run_pending()
    time.sleep(1)
