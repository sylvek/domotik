#!/bin/bash
/home/pi/domotik/mosquitto_sub_sensors_watt_cc128.sh &&
/home/pi/domotik/mosquitto_sub_sensors_temp_cc128.sh &&
sleep 15 &&
/home/pi/domotik/mosquitto_sub_sensors_watt_cc128.sh &&
/home/pi/domotik/mosquitto_sub_sensors_temp_cc128.sh &&
sleep 15 &&
/home/pi/domotik/mosquitto_sub_sensors_watt_cc128.sh &&
/home/pi/domotik/mosquitto_sub_sensors_temp_cc128.sh &&
sleep 15 &&
/home/pi/domotik/mosquitto_sub_sensors_watt_cc128.sh &&
/home/pi/domotik/mosquitto_sub_sensors_temp_cc128.sh
