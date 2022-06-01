FROM golang:1.18-alpine AS builder
COPY domotik /domotik
WORKDIR /domotik
RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static
COPY --from=builder /domotik/datastore /datastore
CMD ["/datastore", "mosquitto", "/database/domotik.db"]
