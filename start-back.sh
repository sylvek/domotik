#!/bin/bash
echo "launching services..."
/home/pi/domotik/services/service-calculateMeanPerHour.py domotik_watt_per_hour sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt &
#/home/pi/domotik/services/service-calculateMeanPerHour.py domotik_motion_per_hour sensors/hcsr505/motion sensors/hcsr505mean/motion measures/livingRoomPerHour/motion &
/home/pi/domotik/services/service-lightMeanPerHour.py domotik_light_per_hour measures/meanPerHour/watt triggers/led/blink &
/home/pi/domotik/services/service-calculateSumPerDay.py domotik_watt_per_day measures/meanPerHour/watt measures/sumPerDay/watt current/sumPerDay/watt &
#/home/pi/domotik/services/service-calculateMeanPerDay.py domotik_mean_outside_temp_per_day sensors/thn132n/temp sensors/thn132n/mean measures/outside/temp &
/home/pi/domotik/services/service-alertValueUp.py domotik_alert_watt_consumption measures/sumPerDay/watt sumPerDay sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
#/home/pi/domotik/services/service-alertDetectInstruction.py domotik_alert_detect_intrusion sensors/hcsr505/event 12 sylvain.maucourt@free.fr smaucourt@gmail.com smtp.free.fr &
#/home/pi/domotik/services/service-alertLowBattery.py domotik_thn132n_low_battery sensors/thn132n/battery triggers/led/blink 30 0000ff &
#/home/pi/domotik/services/service-displayToLCD.py domotik_display_to_lcd sensors/+/temp triggers/lcd/text &
/home/pi/domotik/services/service-discoverHotWaterConsumption.py domotik_water_per_day sensors/cc128/watt measures/tankHotWaterPerDay/min &
#/home/pi/domotik/services/service-triggerFromHTTP.py domotik_trigger_from_http triggers/web/receive &
#/home/pi/domotik/services/service-executeCommand.py domotik_trigger_from_event sensors/hcsr505/event "/home/pi/domotik/mosquitto_pub/jpg_webcam.sh sensors/camera/jpg" &
/home/pi/domotik/services/service-lightAmbiant.py domotik_light_ambiant triggers/led2/update &
#/home/pi/domotik/services/service-persistCurrentValue.py domotik_persist_thn132n sensors/thn132n/temp current/thn132n/temp &
/home/pi/domotik/services/service-persistCurrentValue.py domotik_persist_esp12e sensors/esp12e/temp current/esp12e/temp &
/home/pi/domotik/services/service-persistCurrentValue.py domotik_persist_esp8266 sensors/esp8266/temp current/esp8266/temp &
#/home/pi/domotik/services/service-persistCurrentValue.py domotik_persist_camera sensors/camera/jpg current/camera/jpg &
/home/pi/domotik/services/service-checkSensorsAvailability.py domotik_check_esp12e sensors/esp12e/temp 30 triggers/led/blink e06fb9 &
#/home/pi/domotik/services/service-checkSensorsAvailability.py domotik_check_thn132n sensors/thn132n/temp 30 triggers/led/blink e06fb9 &

echo "launching mosquitto publishing..."
#/home/pi/domotik/mosquitto_pub/hcsr505.py &

echo "launching mosquitto subscribing..."
#/home/pi/domotik/mosquitto_sub/syslog.py &
#/home/pi/domotik/mosquitto_sub/mongodb.py &
#/home/pi/domotik/mosquitto_sub/freebox.sh &
/home/pi/domotik/mosquitto_sub/led.py &
