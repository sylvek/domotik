package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/glebarez/go-sqlite"
	database "github.com/sylvek/domotik/datastore/database"
)

type Log struct {
	kind  string
	name  string
	unit  string
	value float32
}

func main() {
	mqttHost := os.Getenv("MQTT_HOST")
	databasePath := os.Getenv("DB_PATH")
	log.Printf("starting - MQTT_HOST:%s - DB_PATH:%s", mqttHost, databasePath)

	database := database.NewDatabase(databasePath)

	logs := make(chan Log, 10)

	client := connect(logs, fmt.Sprintf("tcp://%s:%d", mqttHost, 1883), []string{"sensors", "measures"})

	go treatLogs(logs, database)

	defer func() {
		database.Close()
		client.Disconnect(250)
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Println("ciao.")
}

func treatLogs(logs chan Log, database database.Database) {
	for {
		l := <-logs
		err := database.AddSeries(l.kind, l.name, l.unit, l.value)
		if err != nil {
			log.Printf(" - error - %s", err)
			time.Sleep(time.Second)
			logs <- l
		}
	}
}

func connect(logs chan Log, mqttBroker string, kinds []string) (client mqtt.Client) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload()[:])
		elements := strings.Split(msg.Topic(), "/")
		value, _ := strconv.ParseFloat(payload, 32)
		logs <- Log{kind: elements[0], name: elements[1], unit: elements[2], value: float32(value)}
	})
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, kind := range kinds {
		client.Subscribe(kind+"/#", 1, nil)
	}
	return client
}
