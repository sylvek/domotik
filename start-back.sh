#!/bin/bash
echo "launching services..."
$HOME/domotik/services/service-calculateMeanPerHour.py watt_per_hour sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
$HOME/domotik/services/service-calculateMeanPerHour.py motion_per_hour sensors/hcsr505/motion sensors/hcsr505mean/motion measures/livingRoomPerHour/motion &
$HOME/domotik/services/service-lightMeanPerHour.py light_per_hour measures/meanPerHour/watt triggers/led/blink &
$HOME/domotik/services/service-calculateSumPerDay.py watt_per_day measures/meanPerHour/watt measures/sumPerDay/watt &
$HOME/domotik/services/service-alertValueUp.py alert_watt_consumption measures/sumPerDay/watt 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
$HOME/domotik/services/service-alertValueUp.py alert_water_consumption measures/tankHotWaterPerDay/min 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
$HOME/domotik/services/service-alertLowBattery.py thn132n_low_battery sensors/thn132n/battery triggers/led/blink 30 ffff00 &
$HOME/domotik/services/service-displayToLCD.py display_to_lcd sensors/+/temp triggers/lcd/text &
$HOME/domotik/services/service-discoverHotWaterConsumption.py water_per_day sensors/cc128/watt measures/tankHotWaterPerDay/min &
$HOME/domotik/services/service-triggerFromHTTP.py trigger_from_http triggers/web/receive &
$HOME/domotik/services/service-executeCommand.py trigger_from_event sensors/hcsr505/event "$HOME/domotik/mosquitto_pub/jpg_webcam.sh sensors/camera/jpg" &
echo "launching mosquitto publishing..."
$HOME/domotik/mosquitto_pub/cc128.py /dev/cc128 &
$HOME/domotik/mosquitto_pub/thn132n.py /dev/thn132n &
sudo $HOME/domotik/mosquitto_pub/hcsr505.py &
echo "launching mosquitto subscribing..."
$HOME/domotik/mosquitto_sub/syslog.py &
$HOME/domotik/mosquitto_sub/mongodb.py &
$HOME/domotik/mosquitto_sub/freebox.sh &
sudo $HOME/domotik/mosquitto_sub/led.py &
