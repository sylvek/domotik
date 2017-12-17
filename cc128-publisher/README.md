# Back side

## Build it

```
$> docker build -t domotik-cc128-publisher .
```

## Run it

```
$> docker run -d --name mosquitto -p 1883:1883 -p 9001:9001 toke/mosquitto:release-1.4.10-2
$> docker run -d --name domotik-cc128-publisher --link mosquitto:mosquitto --device /dev/cc128:/dev/cc128 domotik-cc128-publisher
```

## Take a snapshot

```
$> docker kill -s SIGUSR1 domotik-back
```
