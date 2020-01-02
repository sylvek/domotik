package com.github.sylvek.domotik.analyzer;

import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonObject;

import java.time.Instant;
import java.util.Calendar;
import java.util.TimeZone;
import java.util.concurrent.atomic.AtomicLong;

public class ConsumptionSumPerDayVerticle extends AbstractConsumptionVerticle {

  private static final String SUM_PER_DAY = "sumPerDay";
  private final AtomicLong sum = new AtomicLong(0L);

  public ConsumptionSumPerDayVerticle() {
    super(EVENT);

    flux()
      .map(jsonObject -> jsonObject.b.getInteger("value"))
      .subscribe(value -> {

        final MessagingService messagingService = MessagingService.eventBus(getVertx());

        long _sum;
        if (countdownTimer()) {
          _sum = sum.getAndSet(value);
          if (_sum > 0)
            messagingService.publish("measures/" + SUM_PER_DAY + "/watt", _sum * 10 / 3_600.0);
        } else {
          _sum = sum.addAndGet(value);
          messagingService.publish("sensors/" + SUM_PER_DAY + "/watt", _sum * 10 / 3_600.0);
        }
      });
  }

  protected long millisecondsForTheNextStep(Instant now) {
    final Calendar c = Calendar.getInstance(TimeZone.getTimeZone("Europe/Paris"));
    c.set(Calendar.HOUR, 23);
    c.set(Calendar.MINUTE, 59);
    c.set(Calendar.SECOND, 59);
    c.set(Calendar.MILLISECOND, 0);

    return c.getTimeInMillis() - now.toEpochMilli();
  }

  @Override
  protected void deserialize(Buffer buffer) {
    final JsonObject jsonObject = new JsonObject(buffer);
    this.sum.set(jsonObject.getLong("sum"));
    this.epoch_for_next_step.set(jsonObject.getLong("epoch_for_next_step"));
  }

  @Override
  protected Buffer serialize() {
    return new JsonObject().put("sum", this.sum.get()).put("epoch_for_next_step", this.epoch_for_next_step.get()).toBuffer();
  }

}
