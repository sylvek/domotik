#!/bin/bash
./measureItRetriever.py tmpr | mosquitto_pub -t sensors/cc128/temp -l
