# Back side

## Build it

```
$> docker build -t domotik-weekly-report .
```

## Run it

```
$> docker run -d --name mongodb mongo:2
$> docker run -d --name domotik-weekly-report --link mongodb:mongodb domotik-weekly-report from@mail.com to@mail.com,to@mail.com smtp.host.com 0.0035 0.148 mongodb
```
