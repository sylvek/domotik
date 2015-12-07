#!/bin/bash
echo "launching web interface..."
cd $HOME/domotik/web
npm install
node_modules/bower/bin/bower install
nohup npm start > /dev/null &
IP=$(ifconfig wlan0 | awk 'sub(/inet addr:/,""){print $1}')
echo "go to http://$IP:3000"
