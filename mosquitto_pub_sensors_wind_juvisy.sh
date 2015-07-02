#!/bin/bash
/home/pi/domotik/yahooWeatherRetriever.py 55863490 wind | mosquitto_pub -t sensors/juvisy/wind -l
