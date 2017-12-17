FROM python:2.7
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "pymongo"]

ENTRYPOINT ["python", "/domotik/mosquitto_sub/mongodb.py", "mosquitto", "1883", "mongodb", "27017"]
