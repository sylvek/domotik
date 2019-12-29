package com.github.sylvek.domotik.analyzer;

import io.vertx.core.eventbus.Message;
import reactor.core.publisher.Flux;

import java.time.Duration;
import java.util.List;
import java.util.stream.Collectors;

public class SensorsToEventVerticle extends DomotikVerticle<String> {

  public static final String SENSORS_LINKY_WATT = "sensors/linky/watt";
  public static final String ACTIVITY_DETECTED = "activityDetected";
  public static final String CONSUMPTION_IS_QUIET = "consumptionIsQuiet";

  public SensorsToEventVerticle(int trigger) {
    super(SENSORS);

    final Flux<List<Integer>> listFlux = flux()
      .buffer(Duration.ofSeconds(10))
      .map(strings -> strings.stream().map(Integer::valueOf).collect(Collectors.toList()));

    listFlux.subscribe(elements -> {
      final long avg = Math.round(elements.stream().mapToInt(Integer::intValue).average().orElse(0d));
      sendEvent((avg > trigger) ? ACTIVITY_DETECTED : CONSUMPTION_IS_QUIET, avg);
    });
  }

  @Override
  public boolean accept(Message<String> message) {
    return message.headers().get(TOPIC).equals(SENSORS_LINKY_WATT);
  }
}
