FROM python:2.7
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "xlrd", "pytz"]

ENTRYPOINT ["./start-back.sh", "mosquitto"]
