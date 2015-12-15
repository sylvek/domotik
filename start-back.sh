#!/bin/bash
echo "launching services..."
$HOME/domotik/services/service-calculateMeanPerHour.py sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
$HOME/domotik/services/service-lightMeanPerHour.py measures/meanPerHour/watt triggers/led/blink &
$HOME/domotik/services/service-calculateSumPerDay.py measures/meanPerHour/watt measures/sumPerDay/watt &
$HOME/domotik/services/service-alertValueUp.py measures/sumPerDay/watt 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
$HOME/domotik/services/service-alertValueUp.py measures/tankHotWaterPerDay/min 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
$HOME/domotik/services/service-alertLowBattery.py sensors/thn132n/battery triggers/led/blink 30 ffff00 &
$HOME/domotik/services/service-displayToLCD.py sensors/+/temp triggers/lcd/text &
$HOME/domotik/services/service-discoverHotWaterConsumption.py sensors/cc128/watt measures/tankHotWaterPerDay/min &
$HOME/domotik/services/service-triggerFromHTTP.py triggers/web/receive &
echo "launching mosquitto publishing..."
$HOME/domotik/mosquitto_pub/cc128.sh &
$HOME/domotik/mosquitto_pub/thn132n.sh &
echo "launching mosquitto subscribing..."
$HOME/domotik/mosquitto_sub/syslog.py &
$HOME/domotik/mosquitto_sub/mongodb.py &
$HOME/domotik/mosquitto_sub/freebox.sh &
sudo $HOME/domotik/mosquitto_sub/led.py &
