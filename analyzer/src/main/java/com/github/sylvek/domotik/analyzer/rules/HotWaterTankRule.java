package com.github.sylvek.domotik.analyzer.rules;

import com.github.sylvek.domotik.analyzer.Rule;
import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.json.JsonObject;

import java.time.LocalTime;
import java.util.concurrent.atomic.AtomicLong;

public class HotWaterTankRule implements Rule {

  public static final double HOT_TANK_WATER_POWER = 2_000;
  private final AtomicLong time = new AtomicLong(0L);

  @Override
  public void process(EventBus eventBus, LocalTime now, JsonObject event) {
    /**
     * Rule 1 :
     *  - Given ~ a detected activity during the "cheapest zone" (meaning between 02:04am and 07:04am or 01:04pm and 04:04pm)
     *  - When  ~ the measured power consumption value is greater than the current hot water tank consumption (~2kW)
     *  - Then  ~ it seems that hot water tank is working
     *
     * Rule 2 :
     *  - Given ~ a worked hot water tank during the "cheapest zone"
     *  - When  ~ the current power consumption value is lower than the hot tank water consumption (ex. 300 vs 2000)
     *          ~
     *  - Then  ~ it seems that hot water tank have finished it work
     */
    final boolean isAfter2am4 = now.isAfter(LocalTime.of(2, 4));
    final boolean isBefore7am4 = now.isBefore(LocalTime.of(7, 4));
    final boolean isAfter1pm4 = now.isAfter(LocalTime.of(13, 4));
    final boolean isBefore4pm4 = now.isBefore(LocalTime.of(16, 4));
    final boolean lowTariff = (isAfter2am4 && isBefore7am4) || (isAfter1pm4 && isBefore4pm4);

    if (!lowTariff) {
      time.set(0L);
      return;
    }

    final Integer value = event.getInteger("value");
    final Long timestamp = event.getLong("timestamp");

    if (value > HOT_TANK_WATER_POWER) {
      time.compareAndSet(0L, timestamp);
    }

    final long triggered = time.get();
    if (value < HOT_TANK_WATER_POWER && triggered > 0) {
      time.set(-1L);
      final long delay = (timestamp - triggered) / 60_000;
      eventBus.publish("measures", delay, new DeliveryOptions().addHeader("topic", "measures/tankHotWaterPerDay/min"));
    }
  }
}
