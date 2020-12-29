package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;
import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

import java.time.LocalTime;

@Rule(name = "detectLowTariffRule", description = "detect when we are in low tariff of electricity", priority = 2)
public class DetectLowTariffRule extends BroadcastableAction {

  public DetectLowTariffRule(DomotikRulesEngine.Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when() {
    var now = LocalTime.now();
    var isBefore7am4 = now.isBefore(LocalTime.of(7, 4));
    var isBefore4pm4 = now.isBefore(LocalTime.of(16, 4));
    var isAfter1pm4 = now.isAfter(LocalTime.of(13, 4));
    var isAfter2am4 = now.isAfter(LocalTime.of(2, 4));

    return isAfter2am4 && isBefore7am4 || isAfter1pm4 && isBefore4pm4;
  }

  @Action
  public void then(Facts facts) {
    facts.put("tariffLow", true);
    facts.put("hotWaterStartedAt", 0d);
  }

}
