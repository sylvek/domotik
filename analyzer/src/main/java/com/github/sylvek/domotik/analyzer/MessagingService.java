package com.github.sylvek.domotik.analyzer;

import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.json.JsonObject;

import java.text.DecimalFormat;

public class MessagingService {

  private static final String EVENT = "event";
  public static final String TRIGGER = "trigger";
  public static final String TOPIC = "topic";
  private static final String RETAIN = "retain";

  private static MessagingService _instance = null;
  private static final DecimalFormat df = new DecimalFormat("#.##");

  private final EventBus eventBus;

  private MessagingService(EventBus eventBus) {
    this.eventBus = eventBus;
  }

  public static MessagingService getInstance(EventBus eventBus) {
    if (_instance == null) {
      _instance = new MessagingService(eventBus);
      return _instance;
    }

    if (!_instance.eventBus.equals(eventBus)) {
      throw new RuntimeException("eventBus instance is different oO");
    }

    return _instance;
  }

  public void publish(String topic, Number value) {
    this.eventBus.publish(TRIGGER,
      df.format(value),
      new DeliveryOptions().addHeader(TOPIC, topic));
  }

  public void publish(String topic, String value) {
    this.eventBus.publish(TRIGGER,
      value,
      new DeliveryOptions().addHeader(TOPIC, topic));
  }

  public void publish(String topic, String value, boolean retain) {
    this.eventBus.publish(TRIGGER,
      value,
      new DeliveryOptions().addHeader(TOPIC, topic).addHeader(RETAIN, Boolean.toString(retain)));
  }

  public void event(String topic, String name, Number value) {
    this.eventBus.publish(EVENT,
      new JsonObject()
        .put("name", name)
        .put("value", value)
        .put("unit", "watt")
        .put("timestamp", System.currentTimeMillis()),
      new DeliveryOptions().addHeader(TOPIC, topic)
    );
  }
}
