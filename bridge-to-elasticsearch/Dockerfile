FROM python:2.7
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "elasticsearch"]

ENTRYPOINT ["python", "/domotik/mosquitto_sub/es.py", "mosquitto", "1883", "elasticsearch", "9200"]
