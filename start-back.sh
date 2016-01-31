#!/bin/bash
echo "launching services..."
nohup $HOME/domotik/services/service-calculateMeanPerHour.py domotik_watt_per_hour sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
nohup $HOME/domotik/services/service-calculateMeanPerHour.py domotik_motion_per_hour sensors/hcsr505/motion sensors/hcsr505mean/motion measures/livingRoomPerHour/motion &
nohup $HOME/domotik/services/service-lightMeanPerHour.py domotik_light_per_hour measures/meanPerHour/watt triggers/led/blink &
nohup $HOME/domotik/services/service-calculateSumPerDay.py domotik_watt_per_day measures/meanPerHour/watt measures/sumPerDay/watt &
nohup $HOME/domotik/services/service-alertValueUp.py domotik_alert_watt_consumption measures/sumPerDay/watt 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
nohup $HOME/domotik/services/service-alertValueUp.py domotik_alert_water_consumption measures/tankHotWaterPerDay/min 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
nohup $HOME/domotik/services/service-alertLowBattery.py domotik_thn132n_low_battery sensors/thn132n/battery triggers/led/blink 30 ffff00 &
nohup $HOME/domotik/services/service-displayToLCD.py domotik_display_to_lcd sensors/+/temp triggers/lcd/text &
nohup $HOME/domotik/services/service-discoverHotWaterConsumption.py domotik_water_per_day sensors/cc128/watt measures/tankHotWaterPerDay/min &
nohup  $HOME/domotik/services/service-triggerFromHTTP.py domotik_trigger_from_http triggers/web/receive &
nohup $HOME/domotik/services/service-executeCommand.py domotik_trigger_from_event sensors/hcsr505/event "$HOME/domotik/mosquitto_pub/jpg_webcam.sh sensors/camera/jpg" &
echo "launching mosquitto publishing..."
nohup $HOME/domotik/mosquitto_pub/cc128.py /dev/cc128 &
nohup $HOME/domotik/mosquitto_pub/thn132n.py /dev/thn132n &
nohup sudo $HOME/domotik/mosquitto_pub/hcsr505.py &
echo "launching mosquitto subscribing..."
nohup $HOME/domotik/mosquitto_sub/syslog.py &
#nohup $HOME/domotik/mosquitto_sub/mongodb.py &
nohup $HOME/domotik/mosquitto_sub/freebox.sh &
nohup sudo $HOME/domotik/mosquitto_sub/led.py &
