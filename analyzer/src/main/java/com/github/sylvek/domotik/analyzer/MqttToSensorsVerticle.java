package com.github.sylvek.domotik.analyzer;

import io.netty.handler.codec.mqtt.MqttQoS;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.Handler;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.Message;
import io.vertx.mqtt.MqttClient;
import io.vertx.mqtt.MqttClientOptions;

public class MqttToSensorsVerticle extends AbstractVerticle {

  public static final int MQTT_PORT = 1883;
  public static final String MQTT_TOPIC = "sensors/#";
  public static final String SENSORS = "sensors";
  public static final String TOPIC = "topic";
  public static final String TRIGGER = "trigger";
  public static final String MEASURES = "measures";
  private static final String RETAIN = "retain";
  private final String mqttServer;
  private final boolean devMode;

  public MqttToSensorsVerticle(String... args) {
    this.mqttServer = args[0];
    this.devMode = (args.length > 1) ? Boolean.valueOf(args[1]) : Boolean.FALSE;
  }

  @Override
  public void start() {
    final MqttClient client = connect();

    final Handler<Message<Object>> topic = message -> client.publish(
      message.headers().get(TOPIC),
      Buffer.buffer(message.body().toString()),
      MqttQoS.AT_LEAST_ONCE, false, false
    );

    final Handler<Message<Object>> devNull = objectMessage -> {};

    getVertx().eventBus().consumer(TRIGGER, (devMode) ? devNull : topic);
    getVertx().eventBus().consumer(MEASURES, (devMode) ? devNull : topic);
  }

  private MqttClient connect() {
    MqttClientOptions options = new MqttClientOptions();
    MqttClient client = MqttClient.create(getVertx(), options);
    client.connect(MQTT_PORT, this.mqttServer, s -> {
      if (s.succeeded()) {
        client.exceptionHandler(h -> {
          h.printStackTrace();
          connect();
        });
        client.subscribe(MQTT_TOPIC, 0);
        client.publishHandler(message -> {
          DeliveryOptions deliveryOptions = new DeliveryOptions().addHeader(TOPIC, message.topicName());
          getVertx().eventBus().publish(SENSORS, message.payload().toString(), deliveryOptions);
        });
      } else {
        System.exit(1);
      }
    });
    return client;
  }
}
