#!/bin/bash
FREEBOX="http://hd1.freebox.fr/pub/remote_control?code=3976946&key="
mosquitto_sub -t triggers/+/menu | xargs curl -s -L "${FREEBOX}power" &
mosquitto_sub -t triggers/+/play | xargs curl -s -L "${FREEBOX}ok" &
mosquitto_sub -t triggers/+/plus | xargs curl -s -L "${FREEBOX}vol_inc" &
mosquitto_sub -t triggers/+/minus | xargs curl -s -L "${FREEBOX}vol_dec" &
mosquitto_sub -t triggers/+/forward | xargs curl -s -L "${FREEBOX}prgm_inc" &
mosquitto_sub -t triggers/+/rewind | xargs curl -s -L "${FREEBOX}prgm_dec" &
