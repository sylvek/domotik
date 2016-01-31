#!/bin/bash
echo "launching services..."
/home/pi/domotik/services/service-calculateMeanPerHour.py domotik_watt_per_hour sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
/home/pi/domotik/services/service-calculateMeanPerHour.py domotik_motion_per_hour sensors/hcsr505/motion sensors/hcsr505mean/motion measures/livingRoomPerHour/motion &
/home/pi/domotik/services/service-lightMeanPerHour.py domotik_light_per_hour measures/meanPerHour/watt triggers/led/blink &
/home/pi/domotik/services/service-calculateSumPerDay.py domotik_watt_per_day measures/meanPerHour/watt measures/sumPerDay/watt &
/home/pi/domotik/services/service-alertValueUp.py domotik_alert_watt_consumption measures/sumPerDay/watt 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
/home/pi/domotik/services/service-alertValueUp.py domotik_alert_water_consumption measures/tankHotWaterPerDay/min 1.20 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
/home/pi/domotik/services/service-alertLowBattery.py domotik_thn132n_low_battery sensors/thn132n/battery triggers/led/blink 30 ffff00 &
/home/pi/domotik/services/service-displayToLCD.py domotik_display_to_lcd sensors/+/temp triggers/lcd/text &
/home/pi/domotik/services/service-discoverHotWaterConsumption.py domotik_water_per_day sensors/cc128/watt measures/tankHotWaterPerDay/min &
/home/pi/domotik/services/service-triggerFromHTTP.py domotik_trigger_from_http triggers/web/receive &
/home/pi/domotik/services/service-executeCommand.py domotik_trigger_from_event sensors/hcsr505/event "/home/pi/domotik/mosquitto_pub/jpg_webcam.sh sensors/camera/jpg" &
echo "launching mosquitto publishing..."
/home/pi/domotik/mosquitto_pub/cc128.py /dev/cc128 &
/home/pi/domotik/mosquitto_pub/thn132n.py /dev/thn132n &
sudo /home/pi/domotik/mosquitto_pub/hcsr505.py &
echo "launching mosquitto subscribing..."
/home/pi/domotik/mosquitto_sub/syslog.py &
/home/pi/domotik/mosquitto_sub/mongodb.py &
/home/pi/domotik/mosquitto_sub/freebox.sh &
/home/pi/domotik/mosquitto_sub/led.py &
