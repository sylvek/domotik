#!/bin/bash
FREEBOX=http://hd1.freebox.fr/pub/remote_control?code=3976946&key=
mosquitto_sub -t sensors/+/menu | xargs -i curl -s -L "${FREEBOX}power" &
mosquitto_sub -t sensors/+/play | xargs -i curl -s -L "${FREEBOX}ok" &
mosquitto_sub -t sensors/+/plus | xargs -i curl -s -L "${FREEBOX}vol_inc" &
mosquitto_sub -t sensors/+/minus | xargs -i curl -s -L "${FREEBOX}vol_dec" &
mosquitto_sub -t sensors/+/forward | xargs -i curl -s -L "${FREEBOX}prgm_inc" &
mosquitto_sub -t sensors/+/rewind | xargs -i curl -s -L "${FREEBOX}prgm_dec" &
