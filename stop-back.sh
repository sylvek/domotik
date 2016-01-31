#!/bin/bash
echo "stopping mosquitto subscribing..."
pkill mosquitto_sub
echo "stopping services..."
pkill domotik
echo "bye"
