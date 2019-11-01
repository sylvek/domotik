package com.github.sylvek.domotik.analyzer.rules;

import io.vertx.core.eventbus.DeliveryOptions;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.json.JsonObject;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;

import java.time.LocalTime;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;

public class HotWaterTankRuleTests {

  private final HotWaterTankRule hotWaterTankRule = new HotWaterTankRule();

  @Test
  @DisplayName("Given no human activity, the hot water tank goes on and off by itself")
  void testWithBasicValues() {
    // Given
    final EventBus eventBus = mock(EventBus.class);
    final LocalTime time = LocalTime.of(2, 05);

    // When
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 300).put("timestamp", 0));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2300).put("timestamp", 60_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2300).put("timestamp", 120_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 300).put("timestamp", 180_000));

    // Then
    ArgumentCaptor<DeliveryOptions> deliveryOptions = ArgumentCaptor.forClass(DeliveryOptions.class);
    ArgumentCaptor<Object> payload = ArgumentCaptor.forClass(Object.class);
    verify(eventBus, times(2)).publish(eq("measures"), payload.capture(), deliveryOptions.capture());

    final List<DeliveryOptions> allDeliveryOptions = deliveryOptions.getAllValues();
    final List<Object> allPayloads = payload.getAllValues();
    assertEquals(2L, allPayloads.get(0));
    assertEquals(1.02, allPayloads.get(1));
    assertEquals("measures/tankHotWaterPerDay/min", allDeliveryOptions.get(0).getHeaders().get("topic"));
    assertEquals("measures/waterPerDay/liter", allDeliveryOptions.get(1).getHeaders().get("topic"));
  }

  @Test
  @DisplayName("Given an human activity during the LT period, the hot water goes on and off by itself")
  void testWithActivities() {
    // Given
    final EventBus eventBus = mock(EventBus.class);
    final LocalTime time = LocalTime.of(2, 05);

    // When
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2300).put("timestamp", 60_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2300).put("timestamp", 120_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 4300).put("timestamp", 180_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 3300).put("timestamp", 240_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 3300).put("timestamp", 300_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 1300).put("timestamp", 360_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 300).put("timestamp", 420_000));

    // Then
    ArgumentCaptor<DeliveryOptions> deliveryOptions = ArgumentCaptor.forClass(DeliveryOptions.class);
    verify(eventBus).publish(eq("measures"), eq(5L), deliveryOptions.capture());

    final List<DeliveryOptions> allValues = deliveryOptions.getAllValues();
    assertEquals("measures/tankHotWaterPerDay/min", allValues.get(0).getHeaders().get("topic"));
  }

  @Test
  @DisplayName("A real case from sunday 13th, october 2019")
  void testRealCase() {
    // Given
    final EventBus eventBus = mock(EventBus.class);
    final LocalTime time = LocalTime.of(2, 05);

    // When
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 276).put("timestamp", 0));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2379).put("timestamp", 60_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2473).put("timestamp", 120_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2361).put("timestamp", 180_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2396).put("timestamp", 240_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2182).put("timestamp", 300_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 253).put("timestamp", 360_000));

    // Then
    ArgumentCaptor<DeliveryOptions> deliveryOptions = ArgumentCaptor.forClass(DeliveryOptions.class);
    verify(eventBus).publish(eq("measures"), eq(5L), deliveryOptions.capture());

    final List<DeliveryOptions> allValues = deliveryOptions.getAllValues();
    assertEquals("measures/tankHotWaterPerDay/min", allValues.get(0).getHeaders().get("topic"));
  }

  @Test
  @DisplayName("Should be done avec a trigger on and off")
  void testShouldBeDone() {
    // Given
    final EventBus eventBus = mock(EventBus.class);
    final LocalTime time = LocalTime.of(2, 05);

    // When
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 276).put("timestamp", 0));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2379).put("timestamp", 60_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2473).put("timestamp", 120_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2361).put("timestamp", 180_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2396).put("timestamp", 240_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2182).put("timestamp", 300_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 253).put("timestamp", 360_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2473).put("timestamp", 420_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2361).put("timestamp", 480_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2396).put("timestamp", 540_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "activityDetected").put("value", 2182).put("timestamp", 600_000));
    hotWaterTankRule.process(eventBus, time, new JsonObject().put("name", "consumptionIsQuiet").put("value", 253).put("timestamp", 660_000));

    // Then
    ArgumentCaptor<DeliveryOptions> deliveryOptions = ArgumentCaptor.forClass(DeliveryOptions.class);
    verify(eventBus, times(2)).publish(eq("measures"), any(), deliveryOptions.capture());

    final List<DeliveryOptions> allValues = deliveryOptions.getAllValues();
    assertEquals("measures/tankHotWaterPerDay/min", allValues.get(0).getHeaders().get("topic"));
    assertEquals("measures/waterPerDay/liter", allValues.get(1).getHeaders().get("topic"));
  }

}
