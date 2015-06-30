#!/bin/bash
/home/pi/domotik/yahooWeatherRetriever.py wind | mosquitto_pub -t sensors/juvisy/wind -l
