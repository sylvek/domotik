#!/bin/bash
mosquitto_sub -v -t sensors/# | /home/pi/domotik/clients/client-graphite.py &
