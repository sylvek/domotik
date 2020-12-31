package com.github.sylvek.domotik;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.file.Path;
import java.util.Calendar;
import java.util.List;
import java.util.Optional;

public class Application {

  private static final Logger LOGGER = LoggerFactory.getLogger(Application.class);

  public static void main(String[] args) {
    var argsList = List.of(args);
    var host = argsList.stream().findFirst().orElse("localhost");
    var prefix = (argsList.contains("--dev")) ? "test/" : "";
    var backup = Path.of(Optional.ofNullable(System.getenv("BACKUP_PATH")).orElse("/domotik/backup.json"));

    LOGGER.info("host: {}", host);
    LOGGER.info("prefix: {}", prefix);
    LOGGER.info("timezone: {}", Calendar.getInstance().getTimeZone().getDisplayName());
    LOGGER.info("backup: {}", backup);

    var domotikService = new DomotikService(host, prefix);
    var domotikRulesEngine = new DomotikRulesEngine(backup, domotikService::publish);

    domotikService.addConsumptionListener(domotikRulesEngine::fire);
    domotikService.start();

    Runtime.getRuntime().addShutdownHook(new Thread(() -> {
      domotikService.stop();
      LOGGER.info("ciao.");
    }));
  }

}
