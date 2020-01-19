# Front side

## Build it

```
$> npm install
$> node_modules/bower/bin/bower install
$> docker build -t domotik-front .
```

## Run it

```
$> docker run -d --name influxdb influx
$> docker run -d --name domotik-front -p 3000:3000 --link influxdb:influxdb domotik-front
```
