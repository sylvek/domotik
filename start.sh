#!/bin/bash
echo "launching mosquitto subscribing..."
cd $HOME/domotik/mosquitto_sub
./mongodb.sh
echo "launching web interface..."
cd $HOME/domotik/web
npm install
node_modules/bower/bin/bower install
nohup npm start &
echo "go to http://localhost:3000"
