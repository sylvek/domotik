# Front side

## Build it

### download wallpapers
```
$> pip3 install image python-resize-image tqdm
$> python3 wallpapers.py
```

### install dependencies
```
$> cd domotik
$> npm install
$> node_modules/bower/bin/bower install
```

### run locally
```
$> npm start
```

## Run it

```
$> docker build -t domotik-front .
$> docker run -d --name influxdb influx
$> docker run -d --name domotik-front -p 3000:3000 --link influxdb:influxdb domotik-front
```
