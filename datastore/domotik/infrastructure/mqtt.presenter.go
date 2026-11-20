package infrastructure

import (
	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/model"
)

// MqttPresenter publishes on the embedded broker.
// It implements port.SummarizePresenter and port.LogPublisher.
type MqttPresenter struct {
	broker *MqttBroker
}

// Present implements port.SummarizePresenter.
func (m *MqttPresenter) Present(output dto.Output) error {
	m.broker.Publish("sensors", "sumPerDay", output.EuroSpentToday, "euro")
	m.broker.Publish("sensors", "sumPerDay", output.RatioLowTariffToday, "rate")
	m.broker.Publish("sensors", "sumPerDay", float64(output.WattConsumedToday), "watt")
	m.broker.Publish("sensors", "meanPerHour", float64(output.WattPerHourForThisHour), "watt")
	m.broker.Publish("sensors", "meanPerMinute", output.WattPerHourForLastMinute, "watt")

	return nil
}

// Publish implements port.LogPublisher.
func (m *MqttPresenter) Publish(l model.Log) error {
	m.broker.Publish(l.GetTopic(), l.GetName(), l.GetValue(), l.GetUnit())
	return nil
}

func NewMqttPresenter(broker *MqttBroker) *MqttPresenter {
	return &MqttPresenter{broker: broker}
}
