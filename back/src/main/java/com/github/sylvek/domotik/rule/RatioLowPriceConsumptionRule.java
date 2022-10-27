package com.github.sylvek.domotik.rule;

import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

import com.github.sylvek.domotik.DomotikRulesEngine.Broadcaster;

@Rule(name = "ratioLowPriceConsumption", description = "calculate the ratio of low price consumption", priority = 5)
public class RatioLowPriceConsumptionRule extends BroadcastableAction {

  public RatioLowPriceConsumptionRule(Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when(@Fact("consumption") double mean) {
    return mean != 0;
  }

  @Action
  public void then(Facts facts) {
    double currentConsumption = facts.get("consumption");
    double lowConsumption = facts.get("currentSumPerDayLowPrice");
    double totalConsumption = facts.get("currentSumPerDay");
    boolean lowPrice = facts.get("lowPrice");

    if (lowPrice) {
      lowConsumption += currentConsumption;
      facts.put("currentSumPerDayLowPrice", lowConsumption);
    }

    this.broadcaster.broadcast("sensors/sumPerDay/rate",
        Double.toString(Math.round(lowConsumption / totalConsumption * 100.0) / 100.0),
        true);
  }
}
