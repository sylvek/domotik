package com.github.sylvek.domotik;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.file.Path;
import java.util.Calendar;
import java.util.List;

public class Application {

  private static final Logger LOGGER = LoggerFactory.getLogger(Application.class);

  public static void main(String[] args) {
    var argsList = List.of(args);
    var host = argsList.stream().findFirst().orElse("localhost");
    var devMode = argsList.contains("--dev");
    var prefix = (devMode) ? "test/" : "";
    var calendar = Calendar.getInstance();
    var backup = Path.of("/domotik/backup.json");

    LOGGER.info("host: " + host);
    LOGGER.info("prefix: " + prefix);
    LOGGER.info("timezone: " + calendar.getTimeZone().getDisplayName());
    LOGGER.info("backup: " + backup);

    var domotikService = new DomotikService(host, prefix);
    var domotikRulesEngine = new DomotikRulesEngine(backup, calendar, domotikService::publish);

    domotikService.addConsumptionListener(domotikRulesEngine::fire);
    domotikService.start();

    Runtime.getRuntime().addShutdownHook(new Thread(() -> {
      domotikService.stop();
      LOGGER.info("ciao.");
    }));
  }

}
