package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

import java.time.LocalTime;

@Rule(name = "detectHighTariffRule", description = "detect when we are in high tariff of electricity", priority = 2)
public class DetectHighTariffRule extends BroadcastableAction {

  public DetectHighTariffRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when(@Fact("tariffLow") boolean isLow) {
    var now = LocalTime.now();
    var isAfter7am4 = now.isAfter(LocalTime.of(7, 4));
    var isAfter4pm4 = now.isAfter(LocalTime.of(16, 4));
    var isBefore1pm4 = now.isBefore(LocalTime.of(13, 4));
    var isBefore2am4 = now.isBefore(LocalTime.of(2, 4));
    var isBeforeMidnight = now.isBefore(LocalTime.of(23, 59));
    var isAfterMidnight = now.isAfter(LocalTime.of(0, 0));

    return isLow && (isAfter7am4 && isBefore1pm4 || isAfter4pm4 && isBeforeMidnight || isAfterMidnight && isBefore2am4);
  }

  @Action
  public void then(Facts facts) {
    facts.put("tariffLow", false);
  }

}
