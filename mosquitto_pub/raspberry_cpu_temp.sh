#!/bin/bash
if [ -z "$1" ]
  then
    echo "No argument supplied"
    exit 1
fi
TOPIC=$1
echo $((`cat /sys/class/thermal/thermal_zone0/temp|cut -c1-2`)).$((`cat /sys/class/thermal/thermal_zone0/temp|cut -c3-5`)) | mosquitto_pub -h 192.168.0.4 -t $TOPIC -s
