# Ambiant Light

## Build it

```
$> docker build -t domotik-ambiant-light .
```

## Run it

```
$> docker run -d --name mosquitto mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
$> docker run -d --name domotik-ambiant-light --link mosquitto:mosquitto domotik-ambiant-light <ip yeelight>
```
