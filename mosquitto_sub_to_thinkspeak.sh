#!/bin/bash
mosquitto_sub -t sensors/cc128/temp | ./client-thinkspeak.py <<api>> <<field>> &
