#!/bin/bash
mosquitto_sub -v -t sensors/# | tr / . | ./client-graphite.py &
