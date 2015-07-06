#!/bin/bash
$HOME/domotik/services/service-measureItRetriever.py http://192.168.0.2/measureit/measureit_public_html tmpr | mosquitto_pub -t sensors/cc128/temp -l
