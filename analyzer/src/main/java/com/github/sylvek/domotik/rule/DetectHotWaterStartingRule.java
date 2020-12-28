package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

@Rule(name = "detectHotWaterStartingRule", description = "detect the beginning of hot water production")
public class DetectHotWaterStartingRule extends BroadcastableAction {

  private static final double HOT_TANK_WATER_POWER = 2_000;

  public DetectHotWaterStartingRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when(@Fact("tariffLow") boolean isLow,
                      @Fact("consumption") double mean,
                      @Fact("hotWaterStartedAt") double started) {
    return isLow && mean > HOT_TANK_WATER_POWER && started == 0d;
  }

  @Action
  public void then(Facts facts) {
    facts.put("hotWaterStartedAt", System.currentTimeMillis());
  }
}
