#!/bin/bash
yesterday=$(date +%s -d "yesterday")
echo "mongoexport --db domotik --collection measures --query '{timestamp:{\$gt:$yesterday}}'" | bash 
