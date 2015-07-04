#!/bin/bash
mosquitto_sub -t sensors/cc128/temp | ./client-thinkspeak.py SMSXCZZDKPNH3HJ5 1 &
mosquitto_sub -t sensors/cc128/watt | ./client-thinkspeak.py SMSXCZZDKPNH3HJ5 2 &
mosquitto_sub -t sensors/juvisy/temp | ./client-thinkspeak.py SMSXCZZDKPNH3HJ5 3 &
mosquitto_sub -t sensors/pi/temp | ./client-thinkspeak.py SMSXCZZDKPNH3HJ5 4 &
mosquitto_sub -t sensors/juvisy/wind | ./client-thinkspeak.py SMSXCZZDKPNH3HJ5 5 &
