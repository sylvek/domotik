package com.github.sylvek.domotik.analyzer.legacy;

import com.github.sylvek.domotik.analyzer.MessagingService;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonObject;

import java.util.Calendar;
import java.util.TimeZone;
import java.util.concurrent.atomic.AtomicLong;

public class ConsumptionMeanPerHourVerticle extends AbstractConsumptionVerticle {

  public static final String MEAN_PER_HOUR = "meanPerHour";
  private final AtomicLong sum = new AtomicLong(0);
  private final AtomicLong nb = new AtomicLong(0);

  public ConsumptionMeanPerHourVerticle() {
    super(EVENT);

    flux()
      .map(jsonObject -> jsonObject.getPayload().getInteger("value"))
      .subscribe(value -> {

        final MessagingService messagingService = MessagingService.eventBus(getVertx());

        long _sum, _nb;
        if (countdownTimer()) {
          _sum = sum.getAndSet(value);
          _nb = nb.getAndSet(1);
          if (_nb > 0 && _sum > 0)
            messagingService.publish("measures/" + MEAN_PER_HOUR + "/watt", _sum / _nb);
        } else {
          _sum = sum.addAndGet(value);
          _nb = nb.addAndGet(1);
          messagingService.publish("sensors/" + MEAN_PER_HOUR + "/watt", _sum / _nb);
        }
      });
  }

  protected long epochForTheNextStep() {
    final Calendar c = Calendar.getInstance(TimeZone.getTimeZone("Europe/Paris"));
    c.add(Calendar.HOUR, 1);
    c.set(Calendar.MINUTE, 0);
    c.set(Calendar.SECOND, 0);
    c.set(Calendar.MILLISECOND, 0);

    return c.getTimeInMillis();
  }

  @Override
  protected void deserialize(Buffer buffer) {
    final JsonObject jsonObject = new JsonObject(buffer);
    this.nb.set(jsonObject.getLong("nb"));
    this.sum.set(jsonObject.getLong("sum"));
    this.epoch_for_next_step.set(jsonObject.getLong("epoch_for_next_step"));
  }

  @Override
  protected Buffer serialize() {
    return new JsonObject().put("sum", this.sum.get()).put("nb", this.nb.get()).put("epoch_for_next_step", this.epoch_for_next_step.get()).toBuffer();
  }

}
