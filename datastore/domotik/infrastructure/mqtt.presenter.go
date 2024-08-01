package infrastructure

import (
	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/port"
)

type MqttPresenter struct {
	client *MqttClient
}

// Present implements port.MqttPresenter.
func (m *MqttPresenter) Present(output dto.Output) error {
	m.client.Publish("sensors", "sumPerDay", output.EuroSpentToday, "euro")
	m.client.Publish("sensors", "sumPerDay", output.RatioLowTariffToday, "rate")
	m.client.Publish("sensors", "sumPerDay", float64(output.WattConsumedToday), "watt")
	m.client.Publish("sensors", "meanPerHour", float64(output.WattPerHourForThisHour), "watt")
	m.client.Publish("sensors", "meanPerMinute", output.WattPerHourForLastMinute, "watt")

	return nil
}

func NewMqttPresenter(client *MqttClient) port.SummarizePresenter {
	return &MqttPresenter{client: client}
}
