package com.github.sylvek.domotik;

import com.hivemq.client.mqtt.MqttGlobalPublishFilter;
import com.hivemq.client.mqtt.mqtt3.Mqtt3BlockingClient;
import com.hivemq.client.mqtt.mqtt3.Mqtt3Client;
import com.hivemq.client.mqtt.mqtt3.reactor.Mqtt3ReactorClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class DomotikService {

  private static final Logger LOGGER = LoggerFactory.getLogger(DomotikService.class);

  private final List<ConsumptionListener> consumptionListeners = new ArrayList<>();
  private final String prefix;
  private final Mqtt3BlockingClient client;

  public DomotikService(String host, String prefix) {
    this.prefix = prefix;
    this.client = Mqtt3Client.builder()
        .identifier(UUID.randomUUID().toString())
        .serverHost(host)
        .buildBlocking();
  }

  public void publish(String topic, String payload) {
    LOGGER.debug("publish: {} -> {}", topic, payload);
    this.client.publishWith().topic(this.prefix + topic).payload(payload.getBytes()).send();
  }

  interface ConsumptionListener {
    void apply(double mean);
  }

  public void addConsumptionListener(ConsumptionListener c) {
    this.consumptionListeners.add(c);
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
        .subscribe(e -> this.consumptionListeners.forEach(action -> action.apply(e.average().orElse(0d))));
    client.subscribeWith().topicFilter("sensors/#").send()
        .getReturnCodes().forEach(s -> LOGGER.info("subscription: {}", s));
  }

  public void stop() {
    client.disconnect();
  }
}
