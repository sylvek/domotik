package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

import java.util.Calendar;

@Rule(name = "detectNewDayRule", description = "detect a new day")
public class DetectNewDayRule extends BroadcastableAction {

  private final Calendar calendar;

  public DetectNewDayRule(Calendar calendar, DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
    this.calendar = calendar;
  }

  @Condition
  public boolean when(@Fact("currentDay") double currentDay) {
    return currentDay != getDay();
  }

  @Action
  public void then(Facts facts) {
    facts.put("currentDay", getDay());
    double currentSumPerDay = facts.get("currentSumPerDay");
    this.broadcaster.broadcast("measures/sumPerDay/watt",
      Long.toString(Math.round(currentSumPerDay * 10 / 3_600.0)));
    facts.put("currentSumPerDay", 0d);
  }

  private double getDay() {
    return calendar.get(Calendar.DAY_OF_YEAR);
  }
}
