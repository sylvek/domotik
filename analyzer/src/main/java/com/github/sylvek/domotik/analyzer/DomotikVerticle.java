package com.github.sylvek.domotik.analyzer;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.eventbus.Message;
import reactor.core.publisher.Flux;
import reactor.core.publisher.FluxSink;

import java.util.ArrayList;
import java.util.List;

public abstract class DomotikVerticle<T> extends AbstractVerticle {

  protected class Tuple<A, B> {
    final A a;
    final B b;

    Tuple(A a, B b) {
      this.a = a;
      this.b = b;
    }
  }

  public static final String SENSORS = "sensors";
  public static final String EVENT = "event";
  public static final String TOPIC = "topic";
  private final List<FluxSink<Tuple<String, T>>> handlers = new ArrayList<>();
  private final Flux<Tuple<String, T>> flux;
  private final String consumerName;

  public DomotikVerticle(String consumerName) {
    super();

    this.consumerName = consumerName;

    this.flux = Flux.create(sink -> {
      handlers.add(sink);
      sink.onCancel(() -> handlers.remove(sink));
    });
  }

  protected Flux<Tuple<String, T>> flux() {
    return flux;
  }

  @Override
  public void start() {
    getVertx().eventBus().consumer(this.consumerName, message -> {
      final Message<T> _m = (Message<T>) message;
      if (accept(_m)) {
        handlers.forEach(handlers -> handlers.next(new Tuple<String, T>(_m.headers().get(TOPIC).split("/")[1], _m.body())));
      }
    });
  }

  public boolean accept(Message<T> message) {
    return true;
  }
}
