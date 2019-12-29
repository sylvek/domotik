package com.github.sylvek.domotik.analyzer;

import io.vertx.core.Vertx;

import java.util.OptionalLong;

public class MainApplication {

  public static final int TRIGGER_ACTIVITY = 500;

  public static void main(String[] args) {
    Vertx vertx = Vertx.vertx();
    vertx.deployVerticle(new EchoVerticle());
    vertx.deployVerticle(new MqttToSensorsVerticle(args));
    vertx.deployVerticle(new SensorsToEventVerticle(TRIGGER_ACTIVITY));
    // vertx.deployVerticle(new SensorsToCurrentVerticle());
    vertx.deployVerticle(new EventToRulesVerticle());
    vertx.deployVerticle(new ConsumptionMeanPerHourVerticle(OptionalLong.empty(), OptionalLong.empty()));
    vertx.deployVerticle(new ConsumptionSumPerDayVerticle(OptionalLong.empty()));
  }
}
