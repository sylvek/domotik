# Back side

## Build it

```
$> docker build -t domotik-bridge-to-mongodb .
```

## Run it

```
$> docker run -d --name mongodb mongo:2
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
$> docker run -d --name domotik-bridge-to-mongodb --link mosquitto:mosquitto --link mongodb:mongodb domotik-bridge-to-mongodb
```
