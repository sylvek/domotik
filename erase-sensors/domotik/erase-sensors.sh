#!/bin/bash
yesterday=$(date +%s -d "2 days ago")
mongo $1/domotik --eval "db.sensors.remove({timestamp:{\$lt:$yesterday}})"
