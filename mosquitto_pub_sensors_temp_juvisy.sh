#!/bin/bash
/home/pi/domotik/yahooWeatherRetriever.py 55863490 temp | mosquitto_pub -t sensors/juvisy/temp -l
