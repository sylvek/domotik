#!/bin/bash
echo "launching services..."
python /domotik/services/service-calculateMeanPerHour.py domotik_watt_per_hour sensors/cc128/watt sensors/cc128mean/watt measures/meanPerHour/watt $1 &
python /domotik/services/service-calculateSumPerDay.py domotik_watt_per_day measures/meanPerHour/watt measures/sumPerDay/watt current/sumPerDay/watt $1 &
python /domotik/services/service-persistCurrentValue.py domotik_persist_esp12e sensors/esp12e/temp current/esp12e/temp $1 &
python /domotik/services/service-persistCurrentValue.py domotik_persist_esp8266 sensors/esp8266/temp current/esp8266/temp $1 &
python /domotik/services/service-checkSensorsAvailability.py domotik_check_esp12e sensors/esp12e/temp 30 triggers/led/blink e06fb9 $1 &

trap 'exit 0' SIGINT SIGTERM

echo "launched"
while :
do
  sleep 1
done
