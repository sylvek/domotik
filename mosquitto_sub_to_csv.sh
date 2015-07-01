#!/bin/bash
DAY=$(date +"%Y-%m-%d %H:%M:%S")
mosquitto_sub -v -t sensors/# | while IFS= read -r line; do echo "$DAY $line"; done >> export.csv &
