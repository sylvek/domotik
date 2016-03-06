#!/bin/bash
echo "launching X..."
DISPLAY=:0
startx -- -nocursor &
sleep 10
xset -dpms
xset dpms 0 0 0
xset s noblank
