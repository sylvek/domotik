#!/bin/bash
./yahooWeatherRetriever.py | mosquitto_pub -t sensors/juvisy/temp -l
