package com.github.sylvek.domotik.analyzer;

import java.time.Instant;
import java.util.Calendar;
import java.util.OptionalLong;
import java.util.concurrent.atomic.AtomicLong;

public class ConsumptionSumPerDayVerticle extends AbstractConsumptionVerticle {

  private static final String SUM_PER_DAY = "sumPerDay";
  private final AtomicLong sum;

  public ConsumptionSumPerDayVerticle(OptionalLong initialSum) {
    super(EVENT);
    this.sum = new AtomicLong(initialSum.orElse(0L));

    flux()
      .map(jsonObject -> jsonObject.getInteger("value"))
      .subscribe(value -> {

        long _sum, _nb;
        if (countdownTimer()) {
          _sum = sum.getAndSet(value);
          publish("measures", SUM_PER_DAY, _sum * 10 / 3_600.0);
        } else {
          _sum = sum.addAndGet(value);
          publish("sensors", SUM_PER_DAY, _sum * 10 / 3_600.0);
        }
      });
  }

  protected long millisecondsForTheNextStep(Instant now) {
    final Calendar c = Calendar.getInstance();
    c.add(Calendar.DAY_OF_WEEK, 1);
    c.set(Calendar.HOUR, 0);
    c.set(Calendar.MINUTE, 0);
    c.set(Calendar.SECOND, 0);

    return c.getTimeInMillis() - now.toEpochMilli();
  }
}
