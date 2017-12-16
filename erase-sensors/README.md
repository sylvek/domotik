# Back side

## Build it

```
$> docker build -t domotik-erase-sensors .
```

## Run it

```
$> docker run -d --name mongodb mongo:2
$> docker run -d --name domotik-erase-sensors --link mongodb:mongodb domotik-erase-sensors
```
