#!/bin/bash
echo "stopping mosquitto subscribing..."
pkill mosquitto_sub
pkill syslog.py
#pkill mongodb.py
pkill freebox.py
pkill led.py

echo "stopping services..."
pkill service

echo "stopping publishers..."
pkill hcsr505.py

echo "bye"
