#!/bin/bash
$HOME/domotik/services/service-piTempRetriever.sh | mosquitto_pub -t sensors/pi/temp -l
