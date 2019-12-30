package com.github.sylvek.domotik.analyzer;

import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.json.JsonObject;

import java.text.DecimalFormat;
import java.time.Instant;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.util.concurrent.atomic.AtomicLong;

public abstract class AbstractConsumptionVerticle extends DomotikVerticle<JsonObject> {

  protected AtomicLong epoch_for_next_step = new AtomicLong(0L);
  private static DecimalFormat df = new DecimalFormat("#.##");

  public AbstractConsumptionVerticle(String consumerName) {
    super(consumerName);
  }

  protected boolean countdownTimer() {
    final Instant now = ZonedDateTime.now(ZoneId.of("Europe/Paris")).toInstant();
    if (Instant.ofEpochMilli(this.epoch_for_next_step.get()).isBefore(now)) {
      this.epoch_for_next_step.set(now.toEpochMilli() + millisecondsForTheNextStep(now));
      return true;
    }
    return false;
  }

  protected void publish(String type, String name, Number value) {
    getVertx().eventBus().publish("trigger", df.format(value), new DeliveryOptions().addHeader("topic", type + "/" + name + "/watt"));
  }

  protected abstract long millisecondsForTheNextStep(Instant now);
}
