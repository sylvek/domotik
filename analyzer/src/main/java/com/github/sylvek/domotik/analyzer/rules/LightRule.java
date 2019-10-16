package com.github.sylvek.domotik.analyzer.rules;

import com.github.sylvek.domotik.analyzer.Rule;
import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.json.JsonObject;

import java.time.LocalTime;

public class LightRule implements Rule {

  private static final String GREEN = "#009900";
  private static final String YELLOW = "#ffa500";
  private static final String RED = "#ff0000";

  @Override
  public void process(EventBus eventBus, LocalTime now, JsonObject event) {
    /**
     * Rule 1 :
     *  - Given ~ an activity
     *  - When  ~ the measured power consumption value is lower than 1000
     *  - Then  ~ consumption is "normal" => green
     *
     * Rule 2 :
     *  - Given ~ an activity
     *  - When  ~ the measured power consumption value is greater than 1000 but lower than 2000
     *  - Then  ~ consumption is "strange" => yellow
     *
     * Rule 3 :
     *  - Given ~ an activity
     *  - When  ~ the measured power consumption value is greater than 2000
     *  - Then  ~ consumption is "abnormal" => red
     */
    final int value = event.getInteger("value");
    eventBus.publish("trigger", valueToColor(value), new DeliveryOptions().addHeader("topic", "triggers/led/update"));
  }

  private String valueToColor(int value) {
    String result = GREEN;

    if (value > 1000) {
      result = YELLOW;
    }

    if (value > 2000) {
      result = RED;
    }

    return result;
  }
}
