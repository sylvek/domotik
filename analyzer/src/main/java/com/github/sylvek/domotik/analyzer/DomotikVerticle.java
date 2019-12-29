package com.github.sylvek.domotik.analyzer;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.Message;
import io.vertx.core.json.JsonObject;
import reactor.core.publisher.Flux;
import reactor.core.publisher.FluxSink;

import java.util.ArrayList;
import java.util.List;

public abstract class DomotikVerticle<T> extends AbstractVerticle {

  public static final String SENSORS = "sensors";
  public static final String EVENT = "event";
  public static final String TOPIC = "topic";
  private final List<FluxSink<T>> handlers = new ArrayList<>();
  private final Flux<T> flux;
  private final String consumerName;

  public DomotikVerticle(String consumerName) {
    super();

    this.consumerName = consumerName;

    this.flux = Flux.create(sink -> {
      handlers.add(sink);
      sink.onCancel(() -> handlers.remove(sink));
    });
  }

  protected void sendEvent(String name, Number value) {
    getVertx().eventBus().publish(EVENT,
      new JsonObject()
        .put("name", name)
        .put("value", value)
        .put("unit", "watt")
        .put("timestamp", System.currentTimeMillis()),
      new DeliveryOptions().addHeader("topic", EVENT)
    );
  }

  protected Flux<T> flux() {
    return flux;
  }

  @Override
  public void start() {
    getVertx().eventBus().consumer(this.consumerName, message -> {
      final Message<T> _m = (Message<T>) message;
      if (accept(_m)) {
        handlers.forEach(handlers -> handlers.next(_m.body()));
      }
    });
  }

  public boolean accept(Message<T> message) {
    return true;
  }
}
