package com.github.sylvek.domotik.analyzer.rules;

import com.github.sylvek.domotik.analyzer.EventToRulesVerticle;
import com.github.sylvek.domotik.analyzer.MessagingService;
import io.vertx.core.json.JsonObject;

import java.util.concurrent.atomic.AtomicInteger;

public class LightRule implements EventToRulesVerticle.Rule {

  enum Color {

    GREEN("#009900"),
    YELLOW("#ffa500"),
    RED("#ff0000");

    private String hexadecimal;

    Color(String hexadecimal) {
      this.hexadecimal = hexadecimal;
    }

    static Color of(int value) {
      if (value > 2000) {
        return RED;
      }
      if (value > 1000) {
        return YELLOW;
      }
      return GREEN;
    }
  }

  private static final AtomicInteger lastState = new AtomicInteger(-1);

  @Override
  public void process(MessagingService messagingService, JsonObject event) {
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
    final Color color = Color.of(value);
    if (lastState.get() != color.ordinal()) {
      messagingService.publish("triggers/led/update", color.hexadecimal);
      lastState.set(color.ordinal());
    }
  }
}
