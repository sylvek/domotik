FROM ruby:2.6-alpine
RUN ["gem", "install", "mqtt", "influxdb"]
COPY domotik /domotik
WORKDIR /domotik
CMD ["ruby", "main.rb", "mosquitto", "influxdb"]