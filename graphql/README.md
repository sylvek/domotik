# graphql side

## Build it

```
$> docker build -t domotik-graphql .
```

## Run it

```
$> docker run -d --name mongodb mongo:2
$> docker run -ti --rm --link mongodb:mongodb -p 3000:3000 domotik-graphql mongodb://mongodb/domotik
```
