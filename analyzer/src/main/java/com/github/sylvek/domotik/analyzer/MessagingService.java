package com.github.sylvek.domotik.analyzer;

import io.vertx.core.Vertx;
import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.EventBus;

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

  public static MessagingService eventBus(Vertx vertx) {
    if (_instance == null) {
      _instance = new MessagingService(vertx.eventBus());
      return _instance;
    }

    if (!_instance.eventBus.equals(vertx.eventBus())) {
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
      new DeliveryOptions().addHeader(TOPIC, topic).addHeader(RETAIN, Boolean.toString(Boolean.TRUE)));
  }
}
