#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import smtplib
import json
import os.path
import csv
import datetime

parser = argparse.ArgumentParser(description='send a mail if value increased by X %')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('value', metavar='value', help='value => sumPerDay, tankHotWaterPerDay, meanPerHour or livingRoomPerHour')
parser.add_argument('from_mail', metavar='from_mail', help='mail sender')
parser.add_argument('to_mail', metavar='to_mail', help='mail receiver')
parser.add_argument('mail_server', metavar='mail_server', help='mail hostname server', nargs='?', default="localhost")
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def get_value_from_csv(data, ext):
    script_dir = os.path.dirname(__file__)
    date = datetime.datetime.now()
    month = date.month
    day = date.isoweekday()
    with open(os.path.join(script_dir, 'data') + '/' + data + '.' + ext, 'rb') as csvfile:
        input_data = csv.reader(csvfile)
        for row in input_data:
            if row[0] == str(month) and row[1] == str(day):
                return float(row[2])
        return -1

def send_mail(current_value):
    header = "From: %s\n" % args.from_mail
    header += "To: %s\n" % args.to_mail
    header += "Subject: [domotik] the consumption seems abnormal."
    body = "current value %s should be abnormal. Please check your installation." % str(current_value)
    server = smtplib.SMTP(args.mail_server)
    server.sendmail(args.from_mail, args.to_mail, header + body)

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    current_value = float(msg.payload)
    min_value = get_value_from_csv(args.value, "min")
    max_value = get_value_from_csv(args.value, "max")
    if not min_value < current_value < max_value:
        send_mail(current_value)

def signal_handler(sig, frame):
    if sig is not signal.SIGUSR1:
        print "Ending and cleaning up"
        client.disconnect()

signal.signal(signal.SIGUSR1, signal_handler)
signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
