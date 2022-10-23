FROM balenalib/raspberry-pi-debian:stretch
RUN apt-get update && apt-get install -y wget libfontconfig1
RUN wget https://dl.grafana.com/oss/release/grafana-8.5.14.linux-armv6.tar.gz
RUN tar -zxvf grafana-8.5.14.linux-armv6.tar.gz
WORKDIR /grafana-8.5.14/bin
ENTRYPOINT ["./grafana-server"]
