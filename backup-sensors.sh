#!/bin/bash
yesterday=$(date +%s -d "yesterday")
echo "mongoexport --db domotik --collection sensors --query '{timestamp:{\$gt:$yesterday}}'" | bash 
