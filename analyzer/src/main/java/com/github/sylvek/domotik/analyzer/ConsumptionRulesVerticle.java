package com.github.sylvek.domotik.analyzer;

import com.github.sylvek.domotik.analyzer.rules.AbnormalActivityRule;
import com.github.sylvek.domotik.analyzer.rules.HotWaterTankRule;
import com.github.sylvek.domotik.analyzer.rules.LightRule;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.json.JsonObject;
import reactor.core.publisher.Flux;
import reactor.core.publisher.FluxSink;

import java.time.LocalTime;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.List;

public class ConsumptionRulesVerticle extends AbstractVerticle {

  private final List<FluxSink<JsonObject>> handlers = new ArrayList<>();
  private final List<Rule> rules = new ArrayList<>();

  public ConsumptionRulesVerticle() {
    super();

    rules.add(new HotWaterTankRule());
    rules.add(new LightRule());
    rules.add(new AbnormalActivityRule());

    Flux<JsonObject> flux = Flux.create(sink -> {
      handlers.add(sink);
      sink.onCancel(() -> handlers.remove(sink));
    });

    flux.buffer(5).subscribe(events -> Flux.fromIterable(events)
      .distinctUntilChanged(json -> json.getString("name"))
      .subscribe(e -> rules.forEach((r) -> r.process(getVertx().eventBus(), LocalTime.now(ZoneId.of("Europe/Paris")), e))));
  }

  @Override
  public void start() throws Exception {
    getVertx().eventBus().consumer("event", message -> {
      final JsonObject json = JsonObject.mapFrom(message.body());
      handlers.forEach(handlers -> handlers.next(json));
    });
  }
}
