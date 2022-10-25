package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	broker "github.com/sylvek/domotik/datastore/broker"
	database "github.com/sylvek/domotik/datastore/database"
)

func main() {
	mqttHost := os.Getenv("MQTT_HOST")
	databasePath := os.Getenv("DB_PATH")
	log.Printf("starting - MQTT_HOST:%s - DB_PATH:%s", mqttHost, databasePath)

	sqlite := database.NewSqliteClient(databasePath)
	mqtt := broker.NewMQTTBrokerClient(mqttHost, 1883)

	logs := make(chan broker.Log, 10)

	go mqtt.ConnectAndListen(logs)
	go treatLogs(logs, sqlite)

	defer func() {
		sqlite.Close()
		mqtt.Disconnect()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Println("ciao.")
}

func treatLogs(logs chan broker.Log, database database.Database) {
	for {
		l := <-logs
		err := database.AddSeries(l.Topic, l.Name, l.Unit, l.Value)
		if err != nil {
			log.Printf(" - error - %s", err)
			time.Sleep(time.Second)
			logs <- l
		}
	}
}
