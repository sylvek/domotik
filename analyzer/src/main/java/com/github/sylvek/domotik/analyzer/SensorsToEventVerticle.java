package com.github.sylvek.domotik.analyzer;

import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.Message;
import io.vertx.core.json.JsonObject;
import reactor.core.publisher.Flux;

import java.time.Duration;
import java.time.LocalTime;
import java.time.ZoneId;
import java.util.List;
import java.util.stream.Collectors;

public class SensorsToEventVerticle extends DomotikVerticle<String> {

  public static final String ACTIVITY_DETECTED = "activityDetected";
  public static final String SENSORS_LINKY_WATT = "sensors/linky/watt";
  public static final String CONSUMPTION_IS_QUIET = "consumptionIsQuiet";
  public static final String EVENT_GENERATED = "events/consumption/detection";

  public static final String LOW = "low";
  public static final String HIGH = "high";

  public SensorsToEventVerticle(int trigger) {
    super(SENSORS);

    final Flux<List<String>> listFlux = flux()
      .buffer(Duration.ofSeconds(10))
      .map(strings -> strings.stream().map(tuple -> tuple.getPayload()).collect(Collectors.toList()));

    listFlux.subscribe(elements -> {
      final long avg = Math.round(elements.stream().mapToInt(Integer::parseInt).average().orElse(0d));
      final String name = (avg > trigger) ? ACTIVITY_DETECTED : CONSUMPTION_IS_QUIET;
      event(EVENT_GENERATED, name, avg, getTimeSlots());
    });
  }

  private void event(String topic, String name, Number value, String timeSlots) {
    this.getVertx().eventBus().publish(EVENT,
      new JsonObject()
        .put("name", name)
        .put("value", value)
        .put("time_slots", timeSlots)
        .put("unit", "watt")
        .put("timestamp", System.currentTimeMillis()),
      new DeliveryOptions().addHeader(TOPIC, topic)
    );
  }

  private String getTimeSlots() {
    final LocalTime now = LocalTime.now(ZoneId.of("Europe/Paris"));
    final boolean isBefore7am4 = now.isBefore(LocalTime.of(7, 4));
    final boolean isBefore4pm4 = now.isBefore(LocalTime.of(16, 4));
    final boolean isAfter1pm4 = now.isAfter(LocalTime.of(13, 4));
    final boolean isAfter2am4 = now.isAfter(LocalTime.of(2, 4));

    return (isAfter2am4 && isBefore7am4 || isAfter1pm4 && isBefore4pm4)  ? LOW : HIGH;
  }

  @Override
  public boolean accept(Message message) {
    return message.headers().get(TOPIC).equals(SENSORS_LINKY_WATT);
  }
}
