FROM mongo:2
COPY domotik /domotik

WORKDIR /domotik
ENTRYPOINT ["./backup-measures.sh", "mongodb:27017"]
