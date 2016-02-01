# domotik

![screenshot](extras/screenshot.png)

Domotik is a very simple home-automation software based on micro-services principle.
Each service (services folder) does only one thing. A script (mosquitto_pub) pushes the data on MQTT.
A another process (mosquitto_sub) is in charge of dispatch the data somewhere.
An another service (services folder) could calculate the mean or the max value, etc and push it on MQTT, etc, etc.

## installation
### raspberrypi
assume that you have installed a fresh raspbian…

### from ansible

```
domotik/ansible $> ansible-playbook -b -i raspberrypi playbook.yml [--limit @host] [--tag "tags"]
# where tags is:
# refresh
# update
# stop
# start
# node_stop
# node_start
```

## sensors (mosquitto_pub)
several sensors push data over MQTT
- home int. temperature (via CurrentCost ENVI cc128)
- home int. temperature (via esp8266 + DS18B20)
- home ext. temperature (via thn132n + arduino bridge over usb)
- power consumption (via CurrentCost ENVI cc128)

## triggers (mosquitto_pub)
- remote IR control (via lirc)
- motion sensor (via hc sr505)

## analyzers (mosquitto_sub)
several analyzers are available
- push data to syslog
- push data to mongodb
- control the freebox HD
- RGB LED
- LCD

## services
- mean by hour
- sum per day
- alert consumption
- webcam picture
