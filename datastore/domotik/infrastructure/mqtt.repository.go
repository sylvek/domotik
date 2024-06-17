package infrastructure

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sylvek/domotik/datastore/domain"
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
func (m *MqttClient) Notify(output domain.Output) error {
	m.publish("sensors", "sumPerDay", output.GetEuroSpentToday(), "euro")
	m.publish("sensors", "sumPerDay", output.GetRatioLowTariffToday(), "rate")
	m.publish("sensors", "sumPerDay", output.GetWattConsumedToday(), "watt")
	m.publish("sensors", "meanPerHour", output.GetWattPerHourForThisHour(), "watt")
	m.publish("sensors", "meanPerMinute", output.GetWattPerHourForLastMinute(), "watt")

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
