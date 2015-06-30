#!/bin/bash
/home/pi/domotik/measureItRetriever.py tmpr | mosquitto_pub -t sensors/cc128/temp -l
