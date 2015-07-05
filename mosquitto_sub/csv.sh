#!/bin/bash
mosquitto_sub -v -t sensors/# | while IFS= read -r line; do echo "$(date +"%Y-%m-%d %H:%M:%S") $line"; done >> export.csv &
