FROM python:2.7
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "serial"]

ENTRYPOINT ["python", "/domotik/mosquitto_sub/cc128.py", "/dev/cc128", "mosquitto"]
