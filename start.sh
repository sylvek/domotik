#!/bin/bash
echo "launching services..."
cd $HOME/domotik/services
./service-cc128.py &
echo "launching mosquitto subscribing..."
cd $HOME/domotik/mosquitto_sub
./syslog.sh
./mongodb.sh
./freebox.sh
echo "launching web interface..."
cd $HOME/domotik/web
npm install
node_modules/bower/bin/bower install
nohup npm start > /dev/null &
echo "go to http://localhost:3000"
