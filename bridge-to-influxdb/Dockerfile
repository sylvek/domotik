FROM golang:1.18-alpine AS builder
COPY domotik /domotik
WORKDIR /domotik
RUN go build

FROM alpine
COPY --from=builder /domotik/bridge /bridge
CMD ["/bridge", "mosquitto", "influxdb"]
