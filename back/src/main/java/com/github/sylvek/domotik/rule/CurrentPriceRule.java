package com.github.sylvek.domotik.rule;

import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.api.Facts;

import com.github.sylvek.domotik.Application;
import com.github.sylvek.domotik.DomotikRulesEngine.Broadcaster;

@Rule(name = "currentPrice", description = "calculate current daily price", priority = 4)
public class CurrentPriceRule extends BroadcastableAction {

  public static final double KWH_LOW_PRICE = 0.1457;
  public static final double KWH_HIGH_PRICE = 0.1963;

  public CurrentPriceRule(Broadcaster broadcaster) {
    super(broadcaster);
  }

  @Condition
  public boolean when(@Fact("consumption") double mean) {
    return mean != 0;
  }

  @Action
  public void then(Facts facts) {
    double currentConsumption = facts.get("consumption");
    double currentPrice = facts.get("currentPrice");
    boolean lowPrice = facts.get("lowPrice");

    double currentSlotPrice = (lowPrice)
        ? (currentConsumption * KWH_LOW_PRICE) / Application.TICK_IN_SECONDS * 1_000
        : (currentConsumption * KWH_HIGH_PRICE) / Application.TICK_IN_SECONDS * 1_000;
    double newCurrentPrice = currentPrice + currentSlotPrice;
    facts.put("currentPrice", newCurrentPrice);
    this.broadcaster.broadcast("sensors/sumPerDay/euro",
        Double.toString(Math.round(
            newCurrentPrice * 1_000.0) / 1_000.0),
        true);
  }

}
