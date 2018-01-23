# Mood Light!

## Build it

```
$> docker build -t domotik-mood-light .
```

## Run it

```
$> docker run -d --name mosquitto -p 1883:1883 -p 9001:9001 toke/mosquitto:release-1.4.10-2
$> docker run -d --name domotik-mood-light --link mosquitto:mosquitto domotik-mood-light
```
