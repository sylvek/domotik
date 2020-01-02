package com.github.sylvek.domotik.analyzer;

import io.vertx.core.Vertx;

import java.util.Arrays;

public class MainApplication {

  public static final int TRIGGER_ACTIVITY = 500;

  public static void main(String[] args) {

    Arrays.stream(args).forEach(System.out::println);

    Vertx vertx = Vertx.vertx();
    vertx.deployVerticle(new EchoVerticle());
    vertx.deployVerticle(new MqttToSensorsVerticle(args));
    vertx.deployVerticle(new SensorsToEventVerticle(TRIGGER_ACTIVITY));
    vertx.deployVerticle(new SensorsToCurrentVerticle());
    vertx.deployVerticle(new EventToRulesVerticle());
    vertx.deployVerticle(new ConsumptionMeanPerHourVerticle());
    vertx.deployVerticle(new ConsumptionSumPerDayVerticle());

    Runtime.getRuntime().addShutdownHook(new Thread(() -> vertx.close()));
  }
}
