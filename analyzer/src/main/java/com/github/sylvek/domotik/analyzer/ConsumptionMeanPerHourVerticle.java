package com.github.sylvek.domotik.analyzer;

import java.time.Instant;
import java.util.Calendar;
import java.util.OptionalLong;
import java.util.concurrent.atomic.AtomicLong;

public class ConsumptionMeanPerHourVerticle extends AbstractConsumptionVerticle {

  public static final String MEAN_PER_HOUR = "meanPerHour";
  private final AtomicLong sum;
  private final AtomicLong nb;

  public ConsumptionMeanPerHourVerticle(OptionalLong initialSum, OptionalLong initialNb) {
    super(EVENT);
    this.sum = new AtomicLong(initialSum.orElse(0L));
    this.nb = new AtomicLong(initialNb.orElse(0L));

    flux()
      .map(jsonObject -> jsonObject.getInteger("value"))
      .subscribe(value -> {

        long _sum, _nb;
        if (countdownTimer()) {
          _sum = sum.getAndSet(value);
          _nb = nb.getAndSet(1);
          if (_nb > 0)
            publish("measures", MEAN_PER_HOUR, _sum / _nb);
        } else {
          _sum = sum.addAndGet(value);
          _nb = nb.addAndGet(1);
          publish("sensors", MEAN_PER_HOUR, _sum / _nb);
        }
      });
  }

  protected long millisecondsForTheNextStep(Instant now) {
    final Calendar c = Calendar.getInstance();
    c.add(Calendar.HOUR, 1);
    c.set(Calendar.MINUTE, 0);
    c.set(Calendar.SECOND, 0);

    return c.getTimeInMillis() - now.toEpochMilli();
  }

}
