#!/bin/bash
echo "launching services..."
$HOME/domotik/services/service-calculateMeanPerHour.py sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
$HOME/domotik/services/service-calculateSumPerDay.py measures/meanPerHour/watt measures/sumPerDay/watt &
echo "launching mosquitto publishing..."
$HOME/domotik/mosquitto_pub/cc128.sh &
echo "launching mosquitto subscribing..."
$HOME/domotik/mosquitto_sub/syslog.py &
$HOME/domotik/mosquitto_sub/mongodb.py &
$HOME/domotik/mosquitto_sub/freebox.sh &
echo "launching web interface..."
cd $HOME/domotik/web
npm install
node_modules/bower/bin/bower install
nohup npm start > /dev/null &
IP=$(ifconfig wlan0 | awk 'sub(/inet addr:/,""){print $1}')
echo "go to http://$IP:3000"
