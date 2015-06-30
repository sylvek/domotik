#!/bin/bash
mosquitto_sub -v -t sensors/# | logger &
