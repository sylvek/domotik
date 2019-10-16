package com.github.sylvek.domotik.analyzer;

import io.vertx.core.eventbus.EventBus;
import io.vertx.core.json.JsonObject;

import java.time.LocalTime;

public interface Rule {

  public void process(EventBus eventBus, LocalTime now, JsonObject event);
}
