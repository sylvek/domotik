FROM ruby:alpine
RUN ["gem", "install", "mqtt"]
COPY domotik /domotik
WORKDIR /domotik
ENTRYPOINT ["ruby", "/domotik/co2-reporter.rb", "mosquitto"]
