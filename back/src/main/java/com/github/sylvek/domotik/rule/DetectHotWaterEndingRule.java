package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

@Rule(name = "detectHotWaterEndingRule", description = "detect the end of hot water production", priority = 3)
public class DetectHotWaterEndingRule extends BroadcastableAction {

  private static final double HOT_TANK_WATER_POWER = 2_000;
  private static final int TEMP_MAX = 60;
  private static final int TEMP_MIN = 10;
  private static final int TANK_CAPACITY = 200;

  public DetectHotWaterEndingRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when(@Fact("consumption") double mean,
      @Fact("hotWaterStartedAt") double started) {
    return started > 0d && mean < HOT_TANK_WATER_POWER;
  }

  @Action
  public void then(Facts facts) {
    double startedAt = facts.get("hotWaterStartedAt");
    var delay = (System.currentTimeMillis() - startedAt) / 60_000;
    this.broadcaster.broadcast("measures/tankHotWaterPerDay/min", Long.toString((long) delay), false);

    var dt = 2.13 * delay * TEMP_MAX / 1_000;
    var temp = TEMP_MAX - dt;
    var consumed = ((TANK_CAPACITY * temp) - 12_000) / (TEMP_MIN - TEMP_MAX);
    this.broadcaster.broadcast("measures/waterPerDay/liter",
        Double.toString(Math.round(consumed * 100) / 100.0), false);

    facts.put("hotWaterStartedAt", -1d);
  }
}
