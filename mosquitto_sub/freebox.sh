#!/bin/bash
FREEBOX="http://hd1.freebox.fr/pub/remote_control?code=3976946&key="
mosquitto_sub -t triggers/+/menu | xargs -I curl -s -L "${FREEBOX}power" &
mosquitto_sub -t triggers/+/play | xargs -I curl -s -L "${FREEBOX}ok" &
mosquitto_sub -t triggers/+/plus | xargs -I curl -s -L "${FREEBOX}vol_inc" &
mosquitto_sub -t triggers/+/minus | xargs -I curl -s -L "${FREEBOX}vol_dec" &
mosquitto_sub -t triggers/+/forward | xargs -I curl -s -L "${FREEBOX}prgm_inc" &
mosquitto_sub -t triggers/+/rewind | xargs -I curl -s -L "${FREEBOX}prgm_dec" &
