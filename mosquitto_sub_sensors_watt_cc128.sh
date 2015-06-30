#!/bin/bash
./measureItRetriever.py watt | mosquitto_pub -t sensors/cc128/watt -l
