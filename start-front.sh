#!/bin/bash
echo "launching mosquitto subscribing..."
/home/pi/domotik/mosquitto_sub/syslog.py $1 &
/home/pi/domotik/mosquitto_sub/mongodb.py $1 &

echo "launching webserver..."
MONGO_DB='127.0.0.1:27017/domotik' forever start -s /home/pi/domotik/web/bin/www
