# domotik-bridge-to-influxdb

## build it

```
$> docker build -t domotik-bridge-to-influxdb .
```

## run it

```
$> docker run -d --name influxdb -p 8086:8086 -v /mnt/disk1/influxdb:/var/lib/influxdb influxdb:alpine
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
$> docker run -d --name domotik-bridge-to-elasticsearch --link mosquitto:mosquitto --link influxdb:influxdb domotik-bridge-to-influxdb
```

