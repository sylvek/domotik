package com.github.sylvek.domotik.analyzer;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.json.JsonObject;
import reactor.core.publisher.Flux;
import reactor.core.publisher.FluxSink;
import reactor.core.publisher.Mono;

import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

public class ConsumptionAggregateVerticle extends AbstractVerticle {

  public static final String SENSORS_LINKY_WATT = "sensors/linky/watt";
  public static final int TRIGGER = 500;

  private final List<FluxSink<Integer>> handlers = new ArrayList<>();

  public ConsumptionAggregateVerticle() {
    super();

    Flux<Integer> flux = Flux.create(sink -> {
      handlers.add(sink);
      sink.onCancel(() -> handlers.remove(sink));
    });

    final Flux<List<Integer>> listFlux = flux.bufferTimeout(10, Duration.ofSeconds(60));

    listFlux.subscribe(integers -> {
      Mono<Integer> meanFlux = Flux.fromIterable(integers)
        .reduce((i1, i2) -> (i1 + i2) / 2);
      meanFlux.filter(integer -> integer > TRIGGER)
        .subscribe(mean -> sendEvent("activityDetected", mean));
      meanFlux.filter(integer -> integer < TRIGGER)
        .subscribe(mean -> sendEvent("consumptionIsQuiet", mean));
    });
  }

  private void sendEvent(String name, Integer value) {
    getVertx().eventBus().publish("event",
      new JsonObject()
        .put("name", name)
        .put("value", value)
        .put("unit", "watt")
        .put("timestamp", System.currentTimeMillis()),
      new DeliveryOptions().addHeader("topic", "event")
    );
  }

  @Override
  public void start() throws Exception {
    getVertx().eventBus().consumer("logs", message -> {
      if (message.headers().get("topic").equals(SENSORS_LINKY_WATT)) {
        final int payload = Integer.parseInt(message.body().toString());
        handlers.forEach(handlers -> handlers.next(payload));
      }
    });
  }
}
