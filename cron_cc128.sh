#!/bin/bash
/home/pi/domotik/mosquitto_pub/sensors_watt_cc128.sh &&
/home/pi/domotik/mosquitto_pub/sensors_temp_cc128.sh &&
sleep 15 &&
/home/pi/domotik/mosquitto_pub/sensors_watt_cc128.sh &&
/home/pi/domotik/mosquitto_pub/sensors_temp_cc128.sh &&
sleep 15 &&
/home/pi/domotik/mosquitto_pub/sensors_watt_cc128.sh &&
/home/pi/domotik/mosquitto_pub/sensors_temp_cc128.sh 
