#!/bin/bash
echo "stopping mosquitto subscribing..."
killall mosquitto_sub
echo "stopping services..."
sudo killall /usr/bin/python
echo "bye"
