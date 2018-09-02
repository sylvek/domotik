# Back side

## Build it

```
$> docker build -t domotik-back .
```

## Run it

```
$> docker run -d --name mosquitto -p 1883:1883 -p 9001:9001 toke/mosquitto:release-1.4.10-2
$> docker run -d --name domotik-back --link mosquitto:mosquitto -v /var/cache:/var/cache domotik-back
```

## Take a snapshot

```
$> docker exec domotik-back pkill -SIGUSR1 python
```
