#!/bin/bash
yesterday=$(date +%s -d "yesterday")
mongo domotik --eval "db.sensors.remove({timestamp:{\$gt:$yesterday}})"
