package com.github.sylvek.domotik.analyzer.rules;

import com.github.sylvek.domotik.analyzer.EventToRulesVerticle;
import com.github.sylvek.domotik.analyzer.MessagingService;
import io.vertx.core.json.JsonObject;

import java.time.LocalTime;
import java.util.Calendar;

public class AbnormalActivityRule implements EventToRulesVerticle.Rule {
  @Override
  public void process(MessagingService messagingService, LocalTime now, JsonObject event) {
    /**
     * Rule 1 :
     *  - Given ~ an activity during the week day between 08:00am to 01:04pm and between 11:30pm to 02:04am
     *  - When  ~ the measured power consumption value is greater than 1000
     *  - Then  ~ it seems there is an abnormal activity
     */
    final int dayOfWeek = Calendar.getInstance().get(Calendar.DAY_OF_WEEK);
    final boolean isWeekDay = dayOfWeek != Calendar.SATURDAY && dayOfWeek != Calendar.SUNDAY;

    final boolean isAfter8am = now.isAfter(LocalTime.of(8, 0));
    final boolean isAfter23pm30 = now.isAfter(LocalTime.of(23, 15));
    final boolean isBefore1pm4 = now.isBefore(LocalTime.of(13, 4));
    final boolean isBefore2am4 = now.isBefore(LocalTime.of(2, 4));

    final String name = event.getString("name");
    final Integer value = event.getInteger("value");

    if (isWeekDay
      && (isBefore2am4 || (isAfter8am && isBefore1pm4) || isAfter23pm30)
      && "activityDetected".equals(name)) {
      messagingService.event("events/consumption/abnormal", "anomalyDetected", value);
    }
  }
}
