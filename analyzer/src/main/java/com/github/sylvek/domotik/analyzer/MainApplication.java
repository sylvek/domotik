package com.github.sylvek.domotik.analyzer;

import io.vertx.core.Vertx;

public class MainApplication {

  public static void main(String[] args) {
    Vertx vertx = Vertx.vertx();
    vertx.deployVerticle(new EchoVerticle());
    vertx.deployVerticle(new MqttToLogsVerticle(args[0]));
    vertx.deployVerticle(new LogsToEventVerticle());
    vertx.deployVerticle(new EventToRulesVerticle());
    // vertx.deployVerticle(new ConsumptionMeanPerHourVerticle());
    // vertx.deployVerticle(new ConsumptionSumPerDayVerticle());
  }
}
