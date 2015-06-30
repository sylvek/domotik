#!/bin/bash
./piTempRetriever.sh | mosquitto_pub -t sensors/pi/temp -l
