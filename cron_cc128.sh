#!/bin/bash
/home/pi/domotik/mosquitto_pub/watt_cc128.sh &&
/home/pi/domotik/mosquitto_pub/temp_cc128.sh &&
sleep 15 &&
/home/pi/domotik/mosquitto_pub/watt_cc128.sh &&
/home/pi/domotik/mosquitto_pub/temp_cc128.sh &&
sleep 15 &&
/home/pi/domotik/mosquitto_pub/watt_cc128.sh &&
/home/pi/domotik/mosquitto_pub/temp_cc128.sh 
