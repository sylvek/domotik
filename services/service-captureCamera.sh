#!/bin/bash
touch /home/pi/domotik/services/fswebcam.conf
fswebcam -c /home/pi/domotik/services/fswebcam.conf $1
