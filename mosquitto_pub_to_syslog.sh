#!/bin/bash
mosquitto_sub -t sensors/# | logger &
