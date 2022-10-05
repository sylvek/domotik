package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

import java.util.Calendar;

@Rule(name = "detectNewHourRule", description = "detect a new hour", priority = 1)
public class DetectNewHourRule extends BroadcastableAction {

  public DetectNewHourRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when(@Fact("currentHour") double currentHour, @Fact("consumption") double mean) {
    return currentHour != getHour() && mean != 0;
  }

  @Action
  public void then(Facts facts) {
    facts.put("currentHour", getHour());
    double currentSumPerHour = facts.get("currentSumPerHour");
    double currentNumberOfStatementPerHour = facts.get("currentNumberOfStatementPerHour");
    this.broadcaster.broadcast("measures/meanPerHour/watt",
        Long.toString(Math.round(currentSumPerHour / currentNumberOfStatementPerHour)), false);
    facts.put("currentSumPerHour", facts.get("consumption"));
    facts.put("currentNumberOfStatementPerHour", 1d);
  }

  private double getHour() {
    return Calendar.getInstance().get(Calendar.HOUR_OF_DAY);
  }
}
