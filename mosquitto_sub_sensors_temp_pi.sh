#!/bin/bash
/home/pi/domotik/piTempRetriever.sh | mosquitto_pub -t sensors/pi/temp -l
