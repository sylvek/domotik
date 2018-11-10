# Ambiant Light

## Build it

```
$> docker build -t domotik-french-co2-reporter .
```

## Run it

```
$> docker run -d --name mosquitto mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
$> docker run -d --name domotik-french-co2-reporter --link mosquitto:mosquitto domotik-french-co2-reporter
```
