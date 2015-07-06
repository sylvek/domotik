#!/bin/bash
$HOME/domotik/services/service-yahooWeatherRetriever.py 55863490 temp | mosquitto_pub -t sensors/juvisy/temp -l
