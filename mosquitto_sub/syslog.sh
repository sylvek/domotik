#!/bin/bash
mosquitto_sub -v -t sensors/# -t measures/# -t triggers/# | logger &
