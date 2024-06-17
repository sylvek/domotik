package infrastructure

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sylvek/domotik/datastore/domain/model"
	"github.com/sylvek/domotik/datastore/port"
)

type MqttClient struct {
	client mqtt.Client
	topics []string
}

// Close implements port.NotificationRepository.
func (m *MqttClient) Close() {
	m.client.Disconnect(250)
}

// Notify implements port.NotificationRepository.
func (m *MqttClient) Notify(output model.Output) error {
	m.publish("sensors", "sumPerDay", output.EuroSpentToday, "euro")
	m.publish("sensors", "sumPerDay", output.RatioLowTariffToday, "rate")
	m.publish("sensors", "sumPerDay", float64(output.WattConsumedToday), "watt")
	m.publish("sensors", "meanPerHour", float64(output.WattPerHourForThisHour), "watt")
	m.publish("sensors", "meanPerMinute", float64(output.WattPerHourForLastMinute), "watt")

	return nil
}

func (m *MqttClient) publish(topic string, name string, value float64, unit string) {
	token := m.client.Publish(
		fmt.Sprintf("%s/%s/%s", topic, name, unit),
		0,
		true,
		fmt.Sprintf("%.3f", value))
	_ = token.Wait()

	if token.Error() != nil {
		log.Println("unable to publish", topic, token.Error())
	}
}

func NewMQTTClient(client mqtt.Client) port.NotificationRepository {
	return &MqttClient{
		client: client,
		topics: []string{"sensors/+/+"}}
}
