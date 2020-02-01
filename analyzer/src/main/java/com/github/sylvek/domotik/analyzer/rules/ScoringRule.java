package com.github.sylvek.domotik.analyzer.rules;

import com.github.sylvek.domotik.analyzer.EventToRulesVerticle;
import com.github.sylvek.domotik.analyzer.MessagingService;
import io.vertx.core.json.JsonObject;

public class ScoringRule implements EventToRulesVerticle.Rule {

  @Override
  public void process(MessagingService messagingService, JsonObject event) {
    final String name = event.getString("name");
    final String timeSlots = event.getString("time_slots");
    final int value = event.getInteger("value");
  }
}
