FROM eclipse-mosquitto:1.5
RUN echo "listener 1883" >> /mosquitto/config/mosquitto.conf
RUN echo "listener 9883" >> /mosquitto/config/mosquitto.conf
RUN echo "protocol websockets" >> /mosquitto/config/mosquitto.conf