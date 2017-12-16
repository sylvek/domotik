# Back side

## Build it

```
$> docker build -t domotik-backup-measures .
```

## Run it

```
$> docker run -d --name mongodb mongo:2
$> docker run -d --name domotik-backup-measures --link mongodb:mongodb domotik-backup-measures
```
