package com.github.sylvek.domotik;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.hivemq.client.mqtt.MqttGlobalPublishFilter;
import com.hivemq.client.mqtt.mqtt3.Mqtt3BlockingClient;
import com.hivemq.client.mqtt.mqtt3.Mqtt3Client;
import com.hivemq.client.mqtt.mqtt3.reactor.Mqtt3ReactorClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.LongStream;

public class DomotikService {

  private static final Logger LOGGER = LoggerFactory.getLogger(DomotikService.class);

  private static final Gson GSON = new GsonBuilder().create();

  private final List<LinkyListener> linkyListeners = new ArrayList<>();
  private final String prefix;
  private final Mqtt3BlockingClient client;

  public DomotikService(String host, String prefix) {
    this.prefix = prefix;
    this.client = Mqtt3Client.builder()
        .identifier(UUID.randomUUID().toString())
        .serverHost(host)
        .automaticReconnectWithDefaultConfig()
        .buildBlocking();
  }

  public void publish(String topic, String payload, boolean retain) {
    LOGGER.debug("publish: {} -> {}", topic, payload);
    this.client.publishWith().topic(this.prefix + topic).payload(payload.getBytes()).retain(retain).send();
  }

  interface LinkyListener {
    void applyConsumption(double mean);

    void applyState(boolean state);
  }

  public void addConsumptionListener(LinkyListener c) {
    this.linkyListeners.add(c);
  }

  public void start(LinkyListener c) {
    this.addConsumptionListener(c);
    start();
  }

  public void start() {
    LOGGER.info("connection: {}", client.connect().getReturnCode());
    var publishes = Mqtt3ReactorClient.from(client).publishes(MqttGlobalPublishFilter.SUBSCRIBED);
    // --- we convert tele-transmitted data from tasmota devices into sensors topic
    publishes
        .filter(p -> p.getTopic().toString().startsWith("tele/"))
        .filter(p -> p.getPayload().isPresent())
        .subscribe(e -> {
          var device = e.getTopic().toString().split("/")[1];
          var AM2301 = (Map) GSON.fromJson(new String(e.getPayloadAsBytes()), Map.class).get("AM2301");
          client.publishWith()
              .topic(prefix + "sensors/" + device + "/temp")
              .payload(Double.toString((double) AM2301.get("Temperature")).getBytes())
              .retain(true).send();
          client.publishWith()
              .topic(prefix + "sensors/" + device + "/humidity")
              .payload(Double.toString((double) AM2301.get("Humidity")).getBytes())
              .retain(true).send();
          client.publishWith()
              .topic(prefix + "sensors/" + device + "/dewpoint")
              .payload(Double.toString((double) AM2301.get("DewPoint")).getBytes())
              .retain(true).send();
        });
    // ---
    publishes
        .filter(p -> p.getTopic().toString().equals("sensors/linky/indice"))
        .filter(p -> p.getPayload().isPresent())
        .buffer(Duration.ofSeconds(Application.TICK_IN_SECONDS))
        .map(e -> e.stream().mapToLong(f -> Long.parseLong(new String(f.getPayloadAsBytes()))))
        .subscribe(
            e -> this.linkyListeners.forEach(action -> action.applyConsumption(wattHourFromIndices(e))));
    publishes
        .filter(p -> p.getTopic().toString().equals("sensors/linky/state"))
        .filter(p -> p.getPayload().isPresent())
        .map(e -> Arrays.equals(e.getPayloadAsBytes(), "0".getBytes()))
        .subscribe(e -> this.linkyListeners.forEach(action -> action.applyState(e)));
    // ---
    if (client.subscribeWith()
        .addSubscription().topicFilter("sensors/linky/indice").applySubscription()
        .addSubscription().topicFilter("sensors/linky/state").applySubscription()
        .addSubscription().topicFilter("tele/+/SENSOR").applySubscription()
        .send().getReturnCodes().stream().allMatch(s -> !s.isError())) {
      LOGGER.info("listening is started");
    } else {
      throw new RuntimeException("unable to listen topics");
    }
  }

  public void stop() {
    client.disconnect();
  }

  private static double wattHourFromIndices(LongStream values) {
    // watt consumed in TICK_IN_SECONDS -> we have to convert it into wh
    var elements = values.toArray();
    Arrays.sort(elements);
    return (elements[elements.length - 1] - elements[0]) * 3600.0 / Application.TICK_IN_SECONDS;
  }
}
