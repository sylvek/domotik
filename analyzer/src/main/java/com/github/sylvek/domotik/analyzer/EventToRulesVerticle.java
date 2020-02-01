package com.github.sylvek.domotik.analyzer;

import com.github.sylvek.domotik.analyzer.rules.HotWaterTankRule;
import com.github.sylvek.domotik.analyzer.rules.LightRule;
import com.github.sylvek.domotik.analyzer.rules.ScoringRule;
import io.vertx.core.json.JsonObject;
import reactor.core.publisher.Flux;

import java.util.ArrayList;
import java.util.List;

public class EventToRulesVerticle extends DomotikVerticle<JsonObject> {

  public interface Rule {
    void process(MessagingService messagingService, JsonObject event);
  }

  private final List<Rule> rules = new ArrayList<>();

  public EventToRulesVerticle() {
    super(EVENT);

    rules.add(new HotWaterTankRule());
    rules.add(new LightRule());
    rules.add(new ScoringRule());

    flux().buffer(5).subscribe(events -> Flux.fromIterable(events)
      .distinctUntilChanged(json -> json.getPayload().getString("name"))
      .subscribe(e -> {
        final MessagingService messagingService = MessagingService.eventBus(getVertx());
        rules.forEach((r) -> r.process(messagingService, e.getPayload()));
      }));
  }
}
