# my own domotic project based on raspberrypi+mqtt

## installation
### raspberrypi
assume that you have installed a fresh raspbian…

### mqtt
sudo apt-get install mosquitto mosquitto-clients python-mosquitto

### node.js
wget http://node-arm.herokuapp.com/node_latest_armhf.deb
sudo dpkg -i node_latest_armhf.deb

### mongodb
wget https://github.com/tjanson/mongodb-armhf-deb/releases/download/v2.1.1-1/mongodb_2.1.1_armhf.deb
sudo dpkg -i mongodb_2.1.1_armhf.deb
sudo /etc/init.d/mongodb start
sudo update-rc.d mongodb defaults

### domotik
cd /home/pi
git clone https://github.com/sylvek/domotik.git
cd /home/pi/domotik
# watch crontab.txt => crontab -e
cd /home/pi/mosquitto_sub
./syslog.sh
./mongodb.sh

## sensors (mosquitto_pub)
several sensors push data over mqtt (read crontab.txt)
- pi temperature
- home int. temperature (via CurrentCost ENVI cc128, and via ws => https://github.com/lalelunet/measureit)
- power consumption (via CurrentCost ENVI cc128, and via ws => https://github.com/lalelunet/measureit)
- home ext. temperature and wind (via yahoo weather webservice)
- via bluetooth LE usb dongle (later?)

## analyzers (mosquitto_sub)
several analyzers are available
- push data to syslog
- push data to csv
- push data to graphite
- push data to thinkspeak
- push data to mongodb
