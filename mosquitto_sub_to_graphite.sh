#!/bin/bash
mosquitto_sub -v -t sensors/# | ./client-graphite.py &
