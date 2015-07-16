#!/bin/bash
mosquitto_sub -v -t sensors/# -t measures/# -t triggers/# | cut -c 1-80 | logger &
