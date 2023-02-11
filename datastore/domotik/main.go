package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sylvek/domotik/datastore/broker"
	"github.com/sylvek/domotik/datastore/compute"
	"github.com/sylvek/domotik/datastore/database"
)

func main() {
	mqttHost := os.Getenv("MQTT_HOST")
	databasePath := os.Getenv("DB_PATH")
	log.Printf("starting - MQTT_HOST:%s - DB_PATH:%s", mqttHost, databasePath)

	sqlite := database.NewSqliteClient(databasePath)
	broker := broker.NewMQTTBrokerClient(fmt.Sprintf("tcp://%s:1883", mqttHost))
	engine := compute.NewRuleEngineClient(databasePath)

	go broker.ConnectAndListen()

	go handleSensorLogs(engine, broker, sqlite)

	defer func() {
		engine.Stop()
		sqlite.Close()
		broker.Disconnect()
		log.Println("ciao.")
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
