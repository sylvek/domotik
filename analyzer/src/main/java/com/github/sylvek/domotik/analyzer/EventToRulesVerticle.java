package com.github.sylvek.domotik.analyzer;

import com.github.sylvek.domotik.analyzer.rules.AbnormalActivityRule;
import com.github.sylvek.domotik.analyzer.rules.HotWaterTankRule;
import com.github.sylvek.domotik.analyzer.rules.LightRule;
import io.vertx.core.json.JsonObject;
import reactor.core.publisher.Flux;

import java.time.LocalTime;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.List;

public class EventToRulesVerticle extends DomotikVerticle<JsonObject> {

  public interface Rule {
    void process(MessagingService messagingService, LocalTime now, JsonObject event);
  }

  private final List<Rule> rules = new ArrayList<>();

  public EventToRulesVerticle() {
    super(EVENT);

    rules.add(new HotWaterTankRule());
    rules.add(new LightRule());
    rules.add(new AbnormalActivityRule());

    flux().buffer(5).subscribe(events -> Flux.fromIterable(events)
      .distinctUntilChanged(json -> json.b.getString("name"))
      .subscribe(e -> {
        final MessagingService messagingService = MessagingService.eventBus(getVertx());
        final LocalTime now = LocalTime.now(ZoneId.of("Europe/Paris"));
        rules.forEach((r) -> r.process(messagingService, now, e.b));
      }));
  }
}
