FROM ruby:alpine3.17
ENV MQTT_HOST=mosquitto
RUN gem install mqtt
COPY main.rb /main.rb
ENTRYPOINT [ "ruby", "/main.rb" ]