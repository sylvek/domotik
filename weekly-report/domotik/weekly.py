#!/usr/bin/python
import signal
import smtplib
import argparse
import time
from pymongo import MongoClient
from datetime import datetime, timedelta
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

parser = argparse.ArgumentParser(description='analyse data from mongodb and generate a mail report')
parser.add_argument('from_mail', metavar='from_mail', help='mail sender')
parser.add_argument('to_mail', metavar='to_mail', help='mail receiver')
parser.add_argument('mail_server', metavar='mail_server', help='mail hostname server', nargs='?', default="localhost")
parser.add_argument('water_price', metavar='water_price', help='water price', nargs='?', default="0.0035")
parser.add_argument('power_price', metavar='power_price', help='power price', nargs='?', default="0.148")
parser.add_argument('mongodb_hostname', metavar='mongodb_hostname', help='hostname of mongodb server', nargs='?', default="0.0.0.0")
parser.add_argument('mongodb_port', metavar='mongodb_port', help='port of mongodb server', nargs='?', default="27017")
args = parser.parse_args()

this_week_datetime = datetime.today() - timedelta(days=7)
this_week = time.mktime(this_week_datetime.timetuple())
power_consumption = 0
water_consumption = 0
power = 0
water = 0

def send_report():
    msg = MIMEMultipart('alternative')
    msg['Subject'] = "[domotik] Votre rapport de progression hebdomadaire"
    msg['From'] = args.from_mail
    msg['To'] = args.to_mail
    text = "%(water_consumption)s has %(water)s euros and %(power_consumption)s has %(power)s euros" % {'water_consumption': water_consumption, 'water': water, 'power_consumption': power_consumption, 'power': power}
    with open('template.html', 'r') as template:
        html = template.read() % {'water_consumption': water_consumption, 'water': water, 'power_consumption': power_consumption, 'power': power}
    part1 = MIMEText(text, 'plain')
    part2 = MIMEText(html, 'html')
    msg.attach(part1)
    msg.attach(part2)
    s = smtplib.SMTP(args.mail_server)
    s.sendmail(args.from_mail, args.to_mail, msg.as_string())
    s.quit()

def signal_handler(sig, frame):
    if sig is not signal.SIGUSR1:
        print "Ending and cleaning up"
        client_mongodb.close()

signal.signal(signal.SIGUSR1, signal_handler)
signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)

client_mongodb = MongoClient(args.mongodb_hostname, int(args.mongodb_port))
for day in client_mongodb.domotik['measures'].find({'sensor': 'sumPerDay', 'type': 'watt', 'timestamp': {'$gt': this_week}}):
    power_consumption += day['value']
for day in client_mongodb.domotik['measures'].find({'sensor': 'waterPerDay', 'type': 'liter', 'timestamp': {'$gt': this_week}}):
    water_consumption += day['value']

power = power_consumption / 1000 * float(args.power_price)
water = water_consumption * float(args.water_price)

send_report()
