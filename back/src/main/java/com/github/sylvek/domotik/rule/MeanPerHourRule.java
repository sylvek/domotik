package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

@Rule(name = "meanPerDay", description = "calculate hourly power consumption", priority = 4)
public class MeanPerHourRule extends BroadcastableAction {

  public MeanPerHourRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when(@Fact("consumption") double mean) {
    return mean != 0;
  }

  @Action
  public void then(Facts facts) {
    double currentConsumption = facts.get("consumption");
    double currentSumPerHour = facts.get("currentSumPerHour");
    double currentNumberOfStatementPerHour = facts.get("currentNumberOfStatementPerHour");
    double newCurrentSumPerHour = currentConsumption + currentSumPerHour;
    double newCurrentNumberOfStatementPerHour = currentNumberOfStatementPerHour + 1;
    facts.put("currentSumPerHour", newCurrentSumPerHour);
    facts.put("currentNumberOfStatementPerHour", newCurrentNumberOfStatementPerHour);
    this.broadcaster.broadcast("sensors/meanPerMinute/watt",
        Long.toString(Math.round(currentConsumption)),
        true);
    this.broadcaster.broadcast("sensors/meanPerHour/watt",
        Long.toString(Math.round(newCurrentSumPerHour / newCurrentNumberOfStatementPerHour)),
        true);
  }
}
