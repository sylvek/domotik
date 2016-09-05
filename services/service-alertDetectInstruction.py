#!/usr/bin/python
import paho.mqtt.client as mqtt
import argparse
import signal
import datetime
import smtplib

parser = argparse.ArgumentParser(description='detect intrusion')
parser.add_argument('service_name', metavar='service_name', help='name of the current service')
parser.add_argument('measure_in', metavar='measure_in', help='measure path given')
parser.add_argument('nb_hours', metavar='nb_hours', help='ex. 12 for 12hours of latency')
parser.add_argument('from_mail', metavar='from_mail', help='mail sender')
parser.add_argument('to_mail', metavar='to_mail', help='mail receiver')
parser.add_argument('mail_server', metavar='mail_server', help='mail hostname server', nargs='?', default="localhost")
parser.add_argument('hostname', metavar='hostname', help='hostname of mqtt server', nargs='?', default="0.0.0.0")
parser.add_argument('port', metavar='port', help='port of mqtt server', nargs='?', default="1883")
args = parser.parse_args()

def send_mail():
    header = "From: %s\n" % args.from_mail
    header += "To: %s\n" % args.to_mail
    header += "Subject: [domotik] possible intrusion detection."
    body = "no motion was recorded for %s hours, and something is appening right now." % args.nb_hours
    server = smtplib.SMTP(args.mail_server)
    server.sendmail(args.from_mail, args.to_mail, header + body)

def on_connect(client, userdata, flags, rc):
    client.subscribe(args.measure_in)

def on_message(client, userdata, msg):
    global locked_datetime
    current_datetime = datetime.datetime.now()
    trigger_datetime = locked_datetime + datetime.timedelta(hours=int(args.nb_hours))
    if msg.payload and trigger_datetime < current_datetime:
        send_mail()
    locked_datetime = current_datetime

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

locked_datetime = datetime.datetime.now()

client.connect(args.hostname, int(args.port), 60)
client.loop_forever()
