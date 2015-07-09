#!/bin/bash
FREEBOX=http://hd1.freebox.fr/pub/remote_control?code=3976946&key=
mosquitto_sub -v -t sensors/+/menu | curl -s -o /dev/null -L "$FREEBOX\power" &
mosquitto_sub -v -t sensors/+/play | curl -s -o /dev/null -L "$FREEBOX\ok" &
mosquitto_sub -v -t sensors/+/plus | curl -s -o /dev/null -L "$FREEBOX\vol_inc" &
mosquitto_sub -v -t sensors/+/minus | curl -s -o /dev/null -L "$FREEBOX\vol_dec" &
mosquitto_sub -v -t sensors/+/forward | curl -s -o /dev/null -L "$FREEBOX\prgm_inc" &
mosquitto_sub -v -t sensors/+/rewind | curl -s -o /dev/null -L "$FREEBOX\prgm_dec" &
