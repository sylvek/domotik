FROM golang:1.18-alpine AS builder
COPY domotik /domotik
WORKDIR /domotik
RUN go build

FROM alpine
COPY --from=builder /domotik/datastore /datastore
CMD ["/datastore", "mosquitto", "/database/domotik.db"]
