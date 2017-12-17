# Front side

## Build it

```
$> docker build -t domotik-front .
```

## Run it

```
$> docker run -d --name mongodb mongo:2
$> docker run -d --name domotik-front -p 3000:3000 --link mongodb:mongodb domotik-front
```
