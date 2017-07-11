#!/bin/bash
echo "launching X..."
startx -- -nocursor &
sleep 10 && xset -display :0 -dpms && xset -display :0 dpms 0 0 0 && xset -display :0 s noblank && xset -display :0 s off && xset -display :0 s reset && xset -display :0 s noexpose
