package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

@Rule(name = "meanPerDay", description = "calculate hourly power consumption")
public class MeanPerHourRule extends BroadcastableAction {

  public MeanPerHourRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when() {
    return true;
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
    this.broadcaster.broadcast("sensors/meanPerHour/watt",
        Long.toString(Math.round(newCurrentSumPerHour / newCurrentNumberOfStatementPerHour)));
  }
}
