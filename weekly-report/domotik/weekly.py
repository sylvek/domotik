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

this_20days_datetime = datetime.today() - timedelta(days=30)
this_20days = time.mktime(this_20days_datetime.timetuple())

power_consumption = 0
water_consumption = 0

power_20days_consumption = 0
water_20days_consumption = 0

power = 0
water = 0

def send_report():
    msg = MIMEMultipart('alternative')
    msg['Subject'] = "[domotik] Votre rapport de progression hebdomadaire"
    msg['From'] = args.from_mail
    msg['To'] = args.to_mail
    text = "%(water)s euros of water and %(power)s euros of electricity" % {'water': water, 'power': power}
    with open('template.html', 'r') as template:
        html = template.read() % {
        'water_consumption': water_consumption,
        'power_consumption': round(power_consumption / 1000, 2),
        'water': water,
        'power': power,
        'water_7days_mean': round(water_consumption / 7, 2),
        'power_7days_mean': round(power_consumption / 7000, 2),
        'water_20days_mean': round(water_20days_consumption / 20, 2),
        'power_20days_mean': round(power_20days_consumption / 20000, 2)
        }
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
for day in client_mongodb.domotik['measures'].find({'sensor': 'sumPerDay', 'type': 'watt', 'timestamp': {'$gt': this_20days}}):
    power_20days_consumption += day['value']
    if(day['timestamp'] > this_week):
        power_consumption += day['value']
for day in client_mongodb.domotik['measures'].find({'sensor': 'waterPerDay', 'type': 'liter', 'timestamp': {'$gt': this_20days}}):
    water_20days_consumption += day['value']
    if(day['timestamp'] > this_week):
        water_consumption += day['value']

power = round(power_consumption / 1000 * float(args.power_price),2)
water = round(water_consumption * float(args.water_price),2)

send_report()
