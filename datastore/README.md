# domotik-datastore

## build it

```
$> docker build -t domotik-datastore .
```

## run it

```
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
$> docker run -d --name datastore --link mosquitto:mosquitto -v /path/to/database:/database domotik-datastore
```

