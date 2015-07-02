#!/bin/bash
mosquitto_sub -v -t sensors/# | tr / . | ./client-graphite.sh 192.168.0.13 2003 &
