# Back side

## Build it

```
$> docker build -t domotik-bridge-to-elasticsearch .
```

## Run it

```
$> docker run -d --name elasticsearch -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" docker.elastic.co/elasticsearch/elasticsearch:6.0.1
$> docker run -d --name kibana -p 5601:5601 docker.elastic.co/kibana/kibana:6.0.1
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto
$> docker run -d --name domotik-bridge-to-elasticsearch --link mosquitto:mosquitto --link elasticsearch:elasticsearch domotik-bridge-to-elasticsearch
```
