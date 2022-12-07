package com.github.sylvek.domotik;

import com.github.sylvek.domotik.rule.*;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.jeasy.rules.api.Facts;
import org.jeasy.rules.api.Rules;
import org.jeasy.rules.api.RulesEngine;
import org.jeasy.rules.core.DefaultRulesEngine;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;

public class DomotikRulesEngine {

  private static final Logger LOGGER = LoggerFactory.getLogger(DomotikRulesEngine.class);

  private final Gson gson;
  private final Path backupPath;

  public interface Broadcaster {
    void broadcast(String topic, String payload, boolean retain);
  }

  private final Facts facts = new Facts();
  private final Rules rules = new Rules();
  private final RulesEngine rulesEngine = new DefaultRulesEngine();

  public DomotikRulesEngine(Path backupPath, Broadcaster broadcaster) {
    this.backupPath = backupPath;
    this.gson = new GsonBuilder().create();

    rules.register(
        new DetectNewDayRule(broadcaster),
        new DetectNewHourRule(broadcaster),
        new SumPerDayRule(broadcaster),
        new MeanPerHourRule(broadcaster),
        new DetectHotWaterStartingRule(broadcaster),
        new DetectHotWaterEndingRule(broadcaster),
        new CurrentPriceRule(broadcaster),
        new RatioLowPriceConsumptionRule(broadcaster));

    restore();
  }

  public void fireConsumption(double mean) {
    facts.put("consumption", mean);
    rulesEngine.fire(rules, facts);
    backup();
  }

  public void fireState(boolean state) {
    facts.put("lowPrice", state);
    if (state) {
      facts.put("hotWaterStartedAt", 0d);
    }
    rulesEngine.fire(rules, facts);
    backup();
  }

  private void backup() {
    String backup = gson.toJson(facts.asMap());
    LOGGER.debug(backup);
    try {
      Files.writeString(this.backupPath, backup);
    } catch (IOException e) {
      LOGGER.error("unable to backup facts");
    }
  }

  private void restore() {
    try {
      var json = Files.readString(this.backupPath);
      Map<String, ?> restore = gson.fromJson(json, Map.class);
      restore.forEach(facts::put);
    } catch (IOException e) {
      LOGGER.warn("unable to restore facts. default value loaded");
      facts.put("consumption", 0d);
      facts.put("lowPrice", false);
      facts.put("currentDay", 0d);
      facts.put("currentHour", 0d);
      facts.put("currentSumPerDay", 0d);
      facts.put("currentSumPerHour", 0d);
      facts.put("currentNumberOfStatementPerHour", 1d);
      facts.put("hotWaterStartedAt", 0d);
      facts.put("currentPrice", 0d);
      facts.put("currentSumPerDayLowPrice", 0d);
    }
  }
}
