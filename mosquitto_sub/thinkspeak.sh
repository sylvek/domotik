#!/bin/bash
mosquitto_sub -t sensors/cc128/temp | ../clients/client-thinkspeak.py <<api>> <<field>> &
