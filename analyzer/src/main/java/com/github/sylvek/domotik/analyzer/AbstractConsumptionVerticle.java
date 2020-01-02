package com.github.sylvek.domotik.analyzer;

import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonObject;

import java.time.Duration;
import java.time.Instant;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.util.concurrent.atomic.AtomicLong;

public abstract class AbstractConsumptionVerticle extends DomotikVerticle<JsonObject> {

  public static final String VAR_CACHE_DOMOTIK_SERVICES = "/var/cache/domotik/services/";
  public static final String EXT = ".json";

  protected AtomicLong epoch_for_next_step = new AtomicLong(0L);

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

  protected abstract long millisecondsForTheNextStep(Instant now);

  protected abstract void deserialize(final Buffer buffer);

  protected abstract Buffer serialize();

  @Override
  public void start() {
    getVertx().fileSystem().mkdirs(VAR_CACHE_DOMOTIK_SERVICES, handler -> {
      final String filename = VAR_CACHE_DOMOTIK_SERVICES + this.getClass().getName() + EXT;
      if (handler.succeeded()) {
        getVertx().fileSystem().readFile(filename, read -> {
          if (read.succeeded()) {
            this.deserialize(read.result());
          }
        });
        flux()
          .buffer(Duration.ofMinutes(1))
          .subscribe(something -> getVertx().fileSystem().writeFile(filename, serialize(), null));
      }
    });
    super.start();
  }
}
