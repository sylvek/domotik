package com.github.sylvek.domotik.analyzer;

import io.netty.handler.codec.mqtt.MqttQoS;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.Handler;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.Message;
import io.vertx.mqtt.MqttClient;
import io.vertx.mqtt.MqttClientOptions;

public class MqttVerticle extends AbstractVerticle {

  public static final int MQTT_PORT = 1883;
  public static final String MQTT_TOPIC = "sensors/#";
  private final String mqttServer;

  public MqttVerticle(String mqttServer) {
    this.mqttServer = mqttServer;
  }

  @Override
  public void start() throws Exception {
    final MqttClient client = connect();
    final Handler<Message<Object>> topic = message -> client.publish(
      message.headers().get("topic"),
      Buffer.buffer(message.body().toString()),
      MqttQoS.AT_LEAST_ONCE, false, false
    );
    getVertx().eventBus().consumer("trigger", topic);
    getVertx().eventBus().consumer("measures", topic);
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
          DeliveryOptions deliveryOptions = new DeliveryOptions();
          deliveryOptions.addHeader("topic", message.topicName());
          getVertx().eventBus().publish("logs", message.payload().toString(), deliveryOptions);
        });
      } else {
        System.exit(1);
      }
    });
    return client;
  }
}
