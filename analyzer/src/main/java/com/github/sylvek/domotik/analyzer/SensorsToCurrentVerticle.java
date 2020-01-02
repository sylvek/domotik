package com.github.sylvek.domotik.analyzer;

import io.vertx.core.eventbus.Message;

import java.util.Arrays;

public class SensorsToCurrentVerticle extends DomotikVerticle<String> {

  public SensorsToCurrentVerticle() {
    super(SENSORS);

    flux().subscribe(tuple -> MessagingService.eventBus(getVertx()).publish("current/" + tuple.a + "/temp", tuple.b, true));
  }

  @Override
  public boolean accept(Message message) {
    return Arrays.asList(new String[]{"sensors/esp32/temp", "sensors/esp8266/temp"}).contains(message.headers().get(TOPIC));
  }
}
