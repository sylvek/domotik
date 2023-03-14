package main

import (
	"log"
	"time"

	"github.com/sylvek/domotik/datastore/broker"
	"github.com/sylvek/domotik/datastore/compute"
	"github.com/sylvek/domotik/datastore/database"
)

func handleSensorLogs(rule compute.Rule, client broker.BrokerClient, database database.Database) {
	for {
		select {
		case output := <-rule.Tick():
			client.Publish(broker.Log{Topic: "sensors", Name: "sumPerDay", Unit: "euro", Value: output.EuroSpentToday})
			client.Publish(broker.Log{Topic: "sensors", Name: "sumPerDay", Unit: "rate", Value: output.RatioLowTariffToday})
			client.Publish(broker.Log{Topic: "sensors", Name: "sumPerDay", Unit: "watt", Value: float64(output.WattConsumedToday)})
			client.Publish(broker.Log{Topic: "sensors", Name: "meanPerHour", Unit: "watt", Value: float64(output.WattPerHourForThisHour)})
			client.Publish(broker.Log{Topic: "sensors", Name: "meanPerMinute", Unit: "watt", Value: float64(output.WattPerHourForLastMinute)})
		case l := <-client.Logs():
			err := database.AddSeries(l.Topic, l.Name, l.Unit, l.Value)
			if err != nil {
				log.Printf(" - error - %s", err)
				time.Sleep(time.Second)
				client.Logs() <- l
			} else {
				if l.Name == "linky" && l.Unit == "indice" {
					rule.SetIndice(int64(l.Value))
				}

				if l.Name == "linky" && l.Unit == "state" {
					rule.SetLowTariffState(l.Value == 0.0)
				}
			}
		}
	}
}
