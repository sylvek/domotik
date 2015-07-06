#!/bin/bash
mosquitto_sub -v -t sensors/# -t measures/# -t triggers/# | $HOME/domotik/clients/client-graphite.py &
