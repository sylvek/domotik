# tasmota-converter

Converts `tele/{your device name}/SENSOR` into a topic based on `sensors/{your device name}/temp`

## build it

```
$> docker build -t tasmota-converter .
```

## run it

```
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
$> docker run -d --name tasmota-converter --link mosquitto:mosquitto tasmota-converter
```
