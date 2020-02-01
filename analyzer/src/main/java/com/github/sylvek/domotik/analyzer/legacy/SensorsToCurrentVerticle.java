package com.github.sylvek.domotik.analyzer.legacy;

import com.github.sylvek.domotik.analyzer.DomotikVerticle;
import com.github.sylvek.domotik.analyzer.MessagingService;
import io.vertx.core.eventbus.Message;

import java.util.Arrays;

public class SensorsToCurrentVerticle extends DomotikVerticle<String> {

  public SensorsToCurrentVerticle() {
    super(SENSORS);

    flux().subscribe(tuple -> MessagingService.eventBus(getVertx()).publish("current/" + tuple.getTopic() + "/temp", tuple.getPayload()));
  }

  @Override
  public boolean accept(Message message) {
    return Arrays.asList(new String[]{"sensors/esp12e/temp"}).contains(message.headers().get(TOPIC));
  }
}
