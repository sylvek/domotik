FROM python:2.7-alpine
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "schedule"]

ENTRYPOINT ["python", "/domotik/services/service-lightAmbiant.py", "service-lightAmbiant", "triggers/led2/update", "mosquitto"]
