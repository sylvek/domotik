# my own domotic project based on raspberrypi+mqtt+graphite

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
