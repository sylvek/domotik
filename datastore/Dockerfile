FROM golang:1.18-alpine AS builder
COPY domotik /domotik
WORKDIR /domotik
RUN CGO_ENABLED=0 go build -o datastore .

FROM gcr.io/distroless/static
COPY --from=builder /domotik/datastore /datastore

ENV SERVER_PORT=3333
ENV MQTT_HOST=tcp://mosquitto:1883
ENV DB_PATH=/database/

EXPOSE 3333

CMD ["/datastore"]
