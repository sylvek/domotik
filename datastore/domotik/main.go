package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sylvek/domotik/datastore/core"
	"github.com/sylvek/domotik/datastore/infrastructure"
)

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	mqttHost := os.Getenv("MQTT_HOST")
	dbPath := os.Getenv("DB_PATH")

	log.Println("datastore starting...")

	client := infrastructure.NewMQTTClient(mqttHost, []string{"sensors/+/+"})
	localRepository := infrastructure.NewLocalStore(dbPath)
	sqliteRepository := infrastructure.NewSqliteDatabase(dbPath)

	application, err := core.NewApplication(localRepository, sqliteRepository)
	if err != nil {
		log.Fatal(err)
	}

	go startMqttListener(client, application)
	go startHttpListener(serverPort, application)

	log.Println("datastore started")

	defer func() {
		client.Stop()
		sqliteRepository.Close()
		localRepository.Close()
		log.Println("datastore stopped")
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func startHttpListener(
	serverPort string,
	application *core.Application) {
	http.HandleFunc("/summarize", func(w http.ResponseWriter, r *http.Request) {
		httpPresenter := infrastructure.NewHttpPresenter(w)
		application.GetSummarize(httpPresenter)
	})

	log.Println("HTTP server started")

	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func startMqttListener(
	client *infrastructure.MqttClient,
	application *core.Application) {

	ticker := time.NewTicker(time.Minute)
	mqttPresenter := infrastructure.NewMqttPresenter(client)

	client.Start()

	log.Println("MQTT listener started")

	for {
		select {
		case <-ticker.C:
			if err := application.Process(mqttPresenter); err != nil {
				log.Fatal(err)
			}
		case msg := <-client.Message():
			payload := string(msg.Payload()[:])
			elements := strings.Split(msg.Topic(), "/")
			value, _ := strconv.ParseFloat(payload, 64)

			err := application.AddLog(elements, value)
			if err != nil {
				log.Printf(" - error - %s", err)
				time.Sleep(time.Second)
				client.Message() <- msg
			}
		}
	}
}
