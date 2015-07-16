# domotik

![screenshot](extras/screenshot.png)

Domotik is a very simple home-automation software based on micro-services principle.
Each service (services folder) does only one thing. A script (mosquitto_pub) pushes the data on MQTT.
A another process (mosquitto_pub) has the responsability to dispatch the data somewhere.
An another service (services folder) could calculate the mean or the max value, etc and push it on MQTT, etc, etc.

## installation
### raspberrypi
assume that you have installed a fresh raspbian…

### mosquitto
- sudo apt-get install mosquitto mosquitto-clients python-mosquitto

- git clone http://git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.python.git
- cd org.eclipse.paho.mqtt.python.git
- sudo python setup.py install

### node.js
- wget http://node-arm.herokuapp.com/node_latest_armhf.deb
- sudo dpkg -i node_latest_armhf.deb

### mongodb
- wget https://github.com/tjanson/mongodb-armhf-deb/releases/download/v2.1.1-1/mongodb_2.1.1_armhf.deb
- sudo dpkg -i mongodb_2.1.1_armhf.deb
- sudo /etc/init.d/mongodb start
- sudo update-rc.d mongodb defaults

- git clone git://github.com/mongodb/mongo-python-driver.git pymongo
- cd pymongo/
- sudo python setup.py install

### others
- sudo apt-get install fswebcam lirc
- http://ozzmaker.com/2013/10/24/how-to-control-the-gpio-on-a-raspberry-pi-with-an-ir-remote/
- http://lirc.sourceforge.net/remotes/apple/A1156

### domotik
- cd /home/pi
- git clone https://github.com/sylvek/domotik.git
- cd /home/pi/domotik
- --> create and run your own sensors (crontab.txt)
- ./start.sh
- go to http://[your raspberrypi]:3000
- ...
- ./stop.sh

## sensors (mosquitto_pub)
several sensors push data over MQTT (read crontab.txt)
- pi temperature
- home int. temperature (via CurrentCost ENVI cc128)
- power consumption (via CurrentCost ENVI cc128)
- home ext. temperature and wind (via yahoo weather webservice)
- webcam usb

## triggers (mosquitto_pub)
- lirc

## analyzers (mosquitto_sub)
several analyzers are available
- push data to syslog
- push data to csv
- push data to graphite
- push data to thinkspeak
- push data to mongodb
- control the freebox HD
