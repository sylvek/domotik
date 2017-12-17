FROM python:2.7
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "pyserial"]

ENTRYPOINT ["python", "/domotik/mosquitto_pub/cc128.py", "/dev/cc128", "mosquitto"]
