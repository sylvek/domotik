#!/bin/bash
mosquitto_sub -v -t sensors/# | ../clients/client-mongodb.py &
