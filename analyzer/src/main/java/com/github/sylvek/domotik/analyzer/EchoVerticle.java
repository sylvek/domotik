package com.github.sylvek.domotik.analyzer;

import io.vertx.core.AbstractVerticle;

import java.time.LocalTime;
import java.time.ZoneId;

public class EchoVerticle extends AbstractVerticle {

  @Override
  public void start() {
    getVertx().eventBus().consumer("event", message -> System.out.println("[" + LocalTime.now(ZoneId.of("Europe/Paris")) + "] " + message.headers().get("topic") + " : " + message.body()));
    getVertx().eventBus().consumer("trigger", message -> System.out.println("[" + LocalTime.now(ZoneId.of("Europe/Paris")) + "] " + message.headers().get("topic") + " : " + message.body()));
    getVertx().eventBus().consumer("measures", message -> System.out.println("[" + LocalTime.now(ZoneId.of("Europe/Paris")) + "] " + message.headers().get("topic") + " : " + message.body()));
  }
}
