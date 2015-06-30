#!/bin/bash
/home/pi/domotik/yahooWeatherRetriever.py | mosquitto_pub -t sensors/juvisy/temp -l
