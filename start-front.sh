#!/bin/bash
echo "launching mosquitto subscribing..."
#/home/pi/domotik/mosquitto_sub/syslog.py $1 &
/home/pi/domotik/mosquitto_sub/mongodb.py $1 &

echo "launching mosquitto publishing..."
/home/pi/domotik/mosquitto_pub/cc128.py /dev/cc128 $1 &
#/home/pi/domotik/mosquitto_pub/thn132n.py /dev/thn132n $1 &

echo "launching services..."
/home/pi/domotik/services/service-rebootIfNecessary.py domotik_raspberry_reboot sensors/cc128/watt 15 $1 &

echo "launching webserver..."
MONGO_DB='127.0.0.1:27017/domotik' forever start -s /home/pi/domotik/web/bin/www
