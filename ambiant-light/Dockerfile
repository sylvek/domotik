FROM ruby:alpine
RUN ["gem", "install", "yeelight-lamp-client", "mqtt", "rx"]
COPY domotik /domotik
WORKDIR /domotik
ENTRYPOINT ["ruby", "/domotik/ambiant-light.rb", "mosquitto"]
