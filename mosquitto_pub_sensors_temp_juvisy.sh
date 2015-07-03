#!/bin/bash
/home/pi/domotik/service_yahooWeatherRetriever.py 55863490 temp | mosquitto_pub -t sensors/juvisy/temp -l
