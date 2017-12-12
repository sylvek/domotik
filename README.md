# domotik

![screenshot](extras/screenshot.png)
![tv](extras/tv.jpg)

Domotik is a very simple home-automation software based on micro-services principle.
Each service (services folder) does only one thing. A script (mosquitto_pub) pushes the data on MQTT.
A another process (mosquitto_sub) is in charge of dispatch the data somewhere.
An another service (services folder) could calculate the mean or the max value, etc and push it on MQTT, etc, etc.

## Build it

```
back> docker build -t domotik-back .
front> docker build -t domotik-front .
bridge-to-mongodb> docker build -t domotik-bridge-to-mongodb .
bridge-to-elasticsearch> docker build -t domotik-bridge-to-elasticsearch .
```

## Run it

```
$> docker run -d --name elasticsearch -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" docker.elastic.co/elasticsearch/elasticsearch:6.0.1
$> docker run -d --name kibana p 5601:5601 docker.elastic.co/kibana/kibana:6.0.1
$> docker run -d --name mongodb mongo:2
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
```
