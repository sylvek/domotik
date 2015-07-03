#!/bin/bash
/home/pi/domotik/service_measureItRetriever.py http://192.168.0.2/measureit/measureit_public_html watt | mosquitto_pub -t sensors/cc128/watt -l
