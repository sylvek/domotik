FROM python:2.7-alpine
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "pytz"]

ENTRYPOINT ["python", "/domotik/services/service-discoverHotWaterConsumption.py", "domotik_water_per_day", "sensors/linky/watt", "measures/tankHotWaterPerDay/min", "0.20", "mosquitto"]
