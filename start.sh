#!/bin/bash
echo "launching services..."
$HOME/domotik/services/service-calculateMeanPerHour.py sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
$HOME/domotik/services/service-calculateSumPerDay.py measures/meanPerHour/watt measures/sumPerDay/watt &
$HOME/domotik/services/service-alertValueUp.py measures/sumPerDay/watt 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
echo "launching mosquitto publishing..."
$HOME/domotik/mosquitto_pub/cc128.sh &
echo "launching mosquitto subscribing..."
$HOME/domotik/mosquitto_sub/syslog.py &
$HOME/domotik/mosquitto_sub/mongodb.py &
$HOME/domotik/mosquitto_sub/freebox.sh &
sudo $HOME/domotik/mosquitto_sub/led.py &
echo "launching web interface..."
cd $HOME/domotik/web
npm install
node_modules/bower/bin/bower install
nohup npm start > /dev/null &
IP=$(ifconfig wlan0 | awk 'sub(/inet addr:/,""){print $1}')
echo "go to http://$IP:3000"
