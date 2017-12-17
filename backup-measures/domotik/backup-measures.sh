#!/bin/bash
yesterday=$(date +%s -d "yesterday")
echo "mongoexport --quiet --host $1 --db domotik --collection measures --query '{timestamp:{\$gt:$yesterday}}'" | bash 
