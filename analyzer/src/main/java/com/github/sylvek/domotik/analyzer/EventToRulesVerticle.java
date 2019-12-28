package com.github.sylvek.domotik.analyzer;

import com.github.sylvek.domotik.analyzer.rules.AbnormalActivityRule;
import com.github.sylvek.domotik.analyzer.rules.HotWaterTankRule;
import com.github.sylvek.domotik.analyzer.rules.LightRule;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.json.JsonObject;
import reactor.core.publisher.Flux;

import java.time.LocalTime;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.List;

public class EventToRulesVerticle extends DomotikVerticle<JsonObject> {

  public interface Rule {
    void process(EventBus eventBus, LocalTime now, JsonObject event);
  }

  private final List<Rule> rules = new ArrayList<>();

  public EventToRulesVerticle() {
    super("event");

    rules.add(new HotWaterTankRule());
    rules.add(new LightRule());
    rules.add(new AbnormalActivityRule());

    flux().buffer(5).subscribe(events -> Flux.fromIterable(events)
      .distinctUntilChanged(json -> json.getString("name"))
      .subscribe(e -> rules.forEach((r) -> r.process(getVertx().eventBus(), LocalTime.now(ZoneId.of("Europe/Paris")), e))));
  }
}
