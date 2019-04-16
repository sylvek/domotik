FROM python:2.7-alpine
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "xlrd", "pytz"]

ENTRYPOINT ["sh", "start-back.sh", "mosquitto"]
