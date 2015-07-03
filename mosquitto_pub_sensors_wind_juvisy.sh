#!/bin/bash
/home/pi/domotik/service_yahooWeatherRetriever.py 55863490 wind | mosquitto_pub -t sensors/juvisy/wind -l
