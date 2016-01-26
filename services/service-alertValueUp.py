#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import smtplib

parser = argparse.ArgumentParser(description='send a mail if value increased by X %')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('percent', metavar='percent', help='10% => 1.10')
parser.add_argument('from_mail', metavar='from_mail', help='mail sender')
parser.add_argument('to_mail', metavar='to_mail', help='mail receiver')
parser.add_argument('mail_server', metavar='mail_server', help='mail hostname server', nargs='?', default="localhost")
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

previous_value = 0

def send_mail(current_value, percent):
    header = "From: %s\n" % args.from_mail
    header += "To: %s\n" % args.to_mail
    header += "Subject: [domotik] consumption increased by %s percent.\n\n" % str(percent)
    body = "current value %s should be abnormal. Please check your installation." % str(current_value)
    server = smtplib.SMTP(args.mail_server)
    server.sendmail(args.from_mail, args.to_mail, header + body)

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    global previous_value
    current_value = float(msg.payload)
    p = float(args.percent)
    if (previous_value != 0 and current_value > previous_value * p):
        send_mail(current_value, current_value * 100 / previous_value)
    previous_value = current_value

def signal_handler(signal, frame):
    print "Ending and cleaning up"
    client.disconnect()

signal.signal(signal.SIGINT, signal_handler)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
