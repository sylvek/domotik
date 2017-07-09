#!/bin/bash
echo "launching X..."
export DISPLAY=:0
startx -- -nocursor &
sleep 10
xset -dpms
xset dpms 0 0 0
xset s noblank
xset s off
xset s reset
