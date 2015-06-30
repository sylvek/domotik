#!/bin/bash
./yahooWeatherRetriever.py wind | mosquitto_pub -t sensors/juvisy/wind -l
