# Domotik

![tv](extras/tv.png)

Domotik is a very simple home-automation software based on micro-services principle.
Each service (services folder) does only one thing. A script (mosquitto_pub) pushes the data on MQTT.
A another process (mosquitto_sub) is in charge of dispatch the data somewhere.
An another service (services folder) could calculate the mean or the max value, etc and push it on MQTT, etc, etc.

## Build it

```
back> docker build -t domotik-back .
front> docker build -t domotik-front .
bridge-to-infludb> docker build -t domotik-bridge-to-influxdb .
```

## Run it

```

$> docker run -d --name influxdb -p 8086:8086 influxdb
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto

$> docker run -d --name domotik-back --link mosquitto:mosquitto domotik-back
$> docker run -d --name domotik-front --link influxdb:influxdb -p 3000:3000 domotik-front

$> docker run -d --name domotik-bridge-to-influxdb --link mosquitto:mosquitto --link influxdb:influxdb domotik-bridge-to-influxdb
```

## Use it

| service | link |
|---------|------|
| tv dashboard | http://ip:3000 |
| MQTT broker | tcp://ip:1883 |
| influxdb | http://ip:8086/domotik |
