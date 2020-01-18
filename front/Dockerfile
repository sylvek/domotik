FROM library/node:13-alpine
COPY domotik /domotik

EXPOSE 3000
WORKDIR /domotik

ENV INFLUX_DB=http://influxdb:8086/domotik

ENTRYPOINT ["npm", "start"]
