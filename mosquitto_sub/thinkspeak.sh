#!/bin/bash
mosquitto_sub -t sensors/cc128/temp | /home/pi/domotik/clients/client-thinkspeak.py <<api>> <<field>> &
