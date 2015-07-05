#!/bin/bash
mosquitto_sub -v -t sensors/# -t measures/# -t triggers/# | /home/pi/domotik/clients/client-graphite.py &
