package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

@Rule(name = "sumPerDay", description = "calculate daily power consumption")
public class SumPerDayRule extends BroadcastableAction {

  public SumPerDayRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when() {
    return true;
  }

  @Action
  public void then(Facts facts) {
    double currentConsumption = facts.get("consumption");
    double currentSumPerDay = facts.get("currentSumPerDay");
    double newCurrentSumPerDay = currentConsumption + currentSumPerDay;
    facts.put("currentSumPerDay", newCurrentSumPerDay);
    this.broadcaster.broadcast(
      "sensors/sumPerDay/watt",
      Long.toString(Math.round(currentSumPerDay * 10 / 3_600.0)));
  }
}
