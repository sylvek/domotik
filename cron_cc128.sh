#!/bin/bash
$HOME/domotik/mosquitto_pub/watt_cc128.sh &&
$HOME/domotik/mosquitto_pub/temp_cc128.sh &&
sleep 15 &&
$HOME/domotik/mosquitto_pub/watt_cc128.sh &&
$HOME/domotik/mosquitto_pub/temp_cc128.sh &&
sleep 15 &&
$HOME/domotik/mosquitto_pub/watt_cc128.sh &&
$HOME/domotik/mosquitto_pub/temp_cc128.sh
