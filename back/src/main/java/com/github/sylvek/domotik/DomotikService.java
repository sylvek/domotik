package com.github.sylvek.domotik;

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
import java.util.UUID;

public class DomotikService {

  private static final Logger LOGGER = LoggerFactory.getLogger(DomotikService.class);

  private final List<LinkyListener> linkyListeners = new ArrayList<>();
  private final String prefix;
  private final Mqtt3BlockingClient client;

  public DomotikService(String host, String prefix) {
    this.prefix = prefix;
    this.client = Mqtt3Client.builder()
        .identifier(UUID.randomUUID().toString())
        .serverHost(host)
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
    publishes
        .filter(p -> p.getTopic().toString().equals("sensors/esp12e/temp"))
        .filter(p -> p.getPayload().isPresent())
        .subscribe(e -> client.publishWith()
            .topic(prefix + "current/esp12e/temp")
            .payload(e.getPayload().orElseThrow())
            .retain(true)
            .send());
    publishes
        .filter(p -> p.getTopic().toString().equals("sensors/linky/watt"))
        .filter(p -> p.getPayload().isPresent())
        .buffer(Duration.ofSeconds(Application.TICK_IN_SECONDS))
        .map(e -> e.stream().mapToInt(f -> Integer.parseInt(new String(f.getPayloadAsBytes()))))
        .subscribe(e -> this.linkyListeners.forEach(action -> action.applyConsumption(e.average().orElse(0d))));
    publishes
        .filter(p -> p.getTopic().toString().equals("sensors/linky/state"))
        .filter(p -> p.getPayload().isPresent())
        .map(e -> Arrays.equals(e.getPayloadAsBytes(), "0".getBytes()))
        .subscribe(e -> this.linkyListeners.forEach(action -> action.applyState(e)));
    client.subscribeWith().topicFilter("sensors/#").send()
        .getReturnCodes().forEach(s -> LOGGER.info("subscription: {}", s));
  }

  public void stop() {
    client.disconnect();
  }
}
