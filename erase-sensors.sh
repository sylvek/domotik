#!/bin/bash
yesterday=$(date +%s -d "2 days ago")
mongo domotik --eval "db.sensors.remove({timestamp:{\$lt:$yesterday}})"
