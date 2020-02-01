package com.github.sylvek.domotik.analyzer;

import com.github.sylvek.domotik.analyzer.legacy.ConsumptionMeanPerHourVerticle;
import com.github.sylvek.domotik.analyzer.legacy.ConsumptionSumPerDayVerticle;
import com.github.sylvek.domotik.analyzer.legacy.SensorsToCurrentVerticle;
import io.vertx.core.Vertx;

import java.util.Arrays;

public class MainApplication {

  public static final int TRIGGER_ACTIVITY = 500;

  public static void main(String[] args) {

    Arrays.stream(args).forEach(System.out::println);

    Vertx vertx = Vertx.vertx();

    /**
     * First of all, {@link MqttToSensorsVerticle} listens sensors/# topic and dispatches messages to eventBus.
     * Also, it allows to publish messages from eventBus (TRIGGER and MEASURE)
     */
    vertx.deployVerticle(new MqttToSensorsVerticle(args));

    /**
     * {@link EchoVerticle} is just a tool verticle that prints messages on terminal
     */
    vertx.deployVerticle(new EchoVerticle());

    /**
     * {@link SensorsToEventVerticle} translate a sensor metric to an event
     */
    vertx.deployVerticle(new SensorsToEventVerticle(TRIGGER_ACTIVITY)); // generates events consumption

    /**
     * This event is processed by rules (thks to {@link EventToRulesVerticle}
     */
    vertx.deployVerticle(new EventToRulesVerticle()); // rules engine

    // -- legacy --
    /**
     * Because when i built my temperature sensors i didn't set the retain value to True,
     * i fixed it by recopying value to an another topic.
     */
    vertx.deployVerticle(new SensorsToCurrentVerticle());

    /**
     * {@link ConsumptionMeanPerHourVerticle} and {@link ConsumptionSumPerDayVerticle} calculate mean and sum
     * in realtime. Very useful.. but Influxdb does the same and could reduce the legacy part.
     */
    vertx.deployVerticle(new ConsumptionMeanPerHourVerticle());
    vertx.deployVerticle(new ConsumptionSumPerDayVerticle());

    Runtime.getRuntime().addShutdownHook(new Thread(() -> vertx.close()));
  }
}
