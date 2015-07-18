#!/bin/bash
echo "launching services..."
$HOME/domotik/services/service-cc128.py &
$HOME/domotik/services/service-calculateMeanPerHour.py sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
$HOME/domotik/services/service-calculateSumPerDay.py measures/meanPerHour/watt measures/sumPerDay/watt &
echo "launching mosquitto subscribing..."
$HOME/domotik/mosquitto_sub/syslog.py &
$HOME/domotik/mosquitto_sub/mongodb.sh
$HOME/domotik/mosquitto_sub/freebox.sh
echo "launching web interface..."
cd $HOME/domotik/web
npm install
node_modules/bower/bin/bower install
nohup npm start > /dev/null &
echo "go to http://localhost:3000"
