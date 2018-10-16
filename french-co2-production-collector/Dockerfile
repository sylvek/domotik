FROM library/alpine

RUN apk add --no-cache curl mosquitto-clients jq
COPY domotik /domotik

WORKDIR /domotik
ENTRYPOINT ["/bin/sh", "run.sh"]
