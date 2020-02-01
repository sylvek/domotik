package com.github.sylvek.domotik.analyzer.rules;

import com.github.sylvek.domotik.analyzer.EventToRulesVerticle;
import com.github.sylvek.domotik.analyzer.MessagingService;
import com.github.sylvek.domotik.analyzer.SensorsToEventVerticle;
import io.vertx.core.json.JsonObject;

import java.util.concurrent.atomic.AtomicLong;

public class HotWaterTankRule implements EventToRulesVerticle.Rule {

  public static final double HOT_TANK_WATER_POWER = 2_000;
  public static final int TEMP_MAX = 60;
  public static final int TEMP_MIN = 10;
  public static final int TANK_CAPACITY = 200;
  private final AtomicLong time = new AtomicLong(0L);

  @Override
  public void process(MessagingService messagingService, JsonObject event) {
    /**
     * Rule 1 :
     *  - Given ~ a detected activity during the "low area" (meaning between 02:04am and 07:04am or 01:04pm and 04:04pm)
     *  - When  ~ the measured power consumption value is greater than the current hot water tank consumption (~2kW)
     *  - Then  ~ it seems that hot water tank is working
     *
     * Rule 2 :
     *  - Given ~ a worked hot water tank during the "low area"
     *  - When  ~ the current power consumption value is lower than the hot tank water consumption (ex. 300 vs 2000)
     *  - Then  ~ it seems that hot water tank have finished it work
     */

    final Integer value = event.getInteger("value");
    final Long timestamp = event.getLong("timestamp");
    final boolean lowTariff = event.getString("time_slots").equals(SensorsToEventVerticle.LOW);

    if (!lowTariff) {
      time.set(0L);
      return;
    }

    if (value > HOT_TANK_WATER_POWER) {
      time.compareAndSet(0L, timestamp);
    }

    final long triggered = time.get();
    if (value < HOT_TANK_WATER_POWER && triggered > 0) {
      time.set(-1L);
      final long delay = (timestamp - triggered) / 60_000;
      messagingService.publish("measures/tankHotWaterPerDay/min", delay);

      final double dt = 2.13 * delay * TEMP_MAX / 1_000;
      final double temp = TEMP_MAX - dt;
      final double consumed = ((TANK_CAPACITY * temp) - 12_000) / (TEMP_MIN - TEMP_MAX);
      messagingService.publish("measures/waterPerDay/liter", Math.round(consumed * 100) / 100.0);
    }
  }
}
