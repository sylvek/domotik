package com.github.sylvek.domotik.analyzer;

import java.time.Instant;
import java.util.Calendar;
import java.util.OptionalLong;
import java.util.TimeZone;
import java.util.concurrent.atomic.AtomicLong;

public class ConsumptionSumPerDayVerticle extends AbstractConsumptionVerticle {

  private static final String SUM_PER_DAY = "sumPerDay";
  private final AtomicLong sum;

  public ConsumptionSumPerDayVerticle(OptionalLong initialSum) {
    super(EVENT);
    this.sum = new AtomicLong(initialSum.orElse(0L));

    flux()
      .map(jsonObject -> jsonObject.b.getInteger("value"))
      .subscribe(value -> {

        final MessagingService messagingService = MessagingService.getInstance(getVertx().eventBus());

        long _sum;
        if (countdownTimer()) {
          _sum = sum.getAndSet(value);
          // if (_sum > 0)
            // messagingService.publish("measures/" + SUM_PER_DAY + "/watt", _sum * 10 / 3_600.0);
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
}
