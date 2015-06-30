#!/bin/bash
(sleep 15 && ./mosquitto_sub_sensors_watt_cc128.sh && ./mosquitto_sub_sensors_temp_cc128.sh) &
(sleep 30 && ./mosquitto_sub_sensors_watt_cc128.sh && ./mosquitto_sub_sensors_temp_cc128.sh) &
(sleep 45 && ./mosquitto_sub_sensors_watt_cc128.sh && ./mosquitto_sub_sensors_temp_cc128.sh) &
(sleep 60 && ./mosquitto_sub_sensors_watt_cc128.sh && ./mosquitto_sub_sensors_temp_cc128.sh) &
