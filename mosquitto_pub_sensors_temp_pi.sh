#!/bin/bash
/home/pi/domotik/service_piTempRetriever.sh | mosquitto_pub -t sensors/pi/temp -l
