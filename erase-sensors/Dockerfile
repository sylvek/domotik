FROM mongo:2
COPY domotik /domotik

WORKDIR /domotik
ENTRYPOINT ["./erase-sensors.sh", "mongodb:27017"]
