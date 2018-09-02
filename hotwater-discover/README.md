# Mood Light!

## Build it

```
$> docker build -t domotik-hotwater-discover .
```

## Run it

```
$> docker run -d --name mosquitto -p 1883:1883 -p 9001:9001 toke/mosquitto:release-1.4.10-2
$> docker run -d --name domotik-hotwater-discover --link mosquitto:mosquitto -v /var/cache:/var/cache domotik-hotwater-discover
```
