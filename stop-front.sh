#!/bin/bash
echo "stopping mosquitto subscribing..."
#pkill syslog.py
pkill mongodb.py
echo "stopping publishers..."
pkill cc128.py
pkill thn132n.py
echo "stopping web server..."
forever stop /home/pi/domotik/web/bin/www
echo "bye"
