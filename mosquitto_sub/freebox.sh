#!/bin/bash
FREEBOX="http://hd1.freebox.fr/pub/remote_control?code=3976946&key="
mosquitto_sub -t triggers/+/menu | xargs -i curl -s -L "${FREEBOX}power" &
mosquitto_sub -t triggers/+/play | xargs -i curl -s -L "${FREEBOX}ok" &
mosquitto_sub -t triggers/+/plus | xargs -i curl -s -L "${FREEBOX}vol_inc" &
mosquitto_sub -t triggers/+/minus | xargs -i curl -s -L "${FREEBOX}vol_dec" &
mosquitto_sub -t triggers/+/forward | xargs -i curl -s -L "${FREEBOX}prgm_inc" &
mosquitto_sub -t triggers/+/rewind | xargs -i curl -s -L "${FREEBOX}prgm_dec" &
