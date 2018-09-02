FROM python:2.7
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "paho-mqtt", "pytz"]

ENTRYPOINT ["python", "/domotik/services/service-discoverHotWaterConsumption.py", "domotik_water_per_day", "sensors/cc128/watt", "measures/tankHotWaterPerDay/min", "0.20", "mosquitto"]
