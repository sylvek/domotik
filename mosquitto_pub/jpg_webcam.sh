#!/bin/bash
if [ -z "$1" ]
  then
    echo "No argument supplied"
    exit 1
fi
TOPIC=$1
CAPTURE=/tmp/capture.jpg
$HOME/domotik/services/service-captureCamera.sh $CAPTURE
base64 -w 0 $CAPTURE | mosquitto_pub -t $TOPIC -s
