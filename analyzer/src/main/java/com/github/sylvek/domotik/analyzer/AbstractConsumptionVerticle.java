package com.github.sylvek.domotik.analyzer;

import io.vertx.core.json.JsonObject;

import java.time.Instant;
import java.util.concurrent.atomic.AtomicLong;

public abstract class AbstractConsumptionVerticle extends DomotikVerticle<JsonObject> {

  protected AtomicLong epoch_for_next_step = new AtomicLong(0L);

  public AbstractConsumptionVerticle(String consumerName) {
    super(consumerName);
  }

  protected boolean countdownTimer() {
    final Instant now = Instant.now();
    if (Instant.ofEpochMilli(this.epoch_for_next_step.get()).isBefore(now)) {
      this.epoch_for_next_step.set(now.toEpochMilli() + millisecondsForTheNextStep(now));
      return true;
    }
    return false;
  }

  protected abstract long millisecondsForTheNextStep(Instant now);
}
