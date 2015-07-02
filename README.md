# my own domotic project based on raspberrypi+mqtt+graphite

## installation
### raspberrypi
assume that you have installed a fresh raspbian…

### mqtt
sudo apt-get install mosquitto mosquitto-clients python-mosquitto

### graphite
https://github.com/dockerana/dockerana

## sensors
several sensors push data over mqtt (read crontab.txt)
- pi temperature
- home int. temperature (via CurrentCost ENVI cc128, and via ws => https://github.com/lalelunet/measureit)
- power consumption (via CurrentCost ENVI cc128, and via ws => https://github.com/lalelunet/measureit)
- home ext. temperature and wind (via yahoo weather webservice)
- via bluetooth LE usb dongle (later?)

## analyzers
several analyzers are available
- push data to syslog
- push data to csv
- push data to graphite
