FROM golang:1.18-alpine AS builder
COPY domotik /domotik
WORKDIR /domotik
RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static
COPY --from=builder /domotik/datastore /datastore

ENV MQTT_HOST=tcp://mosquitto:1883
ENV DB_PATH=/database/

CMD ["/datastore"]
