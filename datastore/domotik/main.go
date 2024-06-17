package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sylvek/domotik/datastore/domain"
	"github.com/sylvek/domotik/datastore/infrastructure"
)

func main() {
	host := os.Getenv("MQTT_HOST")
	path := os.Getenv("DB_PATH")
	log.Printf("starting - MQTT_HOST:%s - DB_PATH:%s", host, path)

	ticker := time.NewTicker(time.Minute)
	client := NewMQTTClient(host, []string{"sensors/+/+"})

	localRepository := infrastructure.NewLocalClient(path)
	mqttRepository := infrastructure.NewMQTTClient(client.GetClient())
	sqliteRepository := infrastructure.NewSqliteDatabase(path)

	defer func() {
		ticker.Stop()
		mqttRepository.Close()
		sqliteRepository.Close()
		localRepository.Close()
		log.Println("ciao.")
	}()

	application, err := NewApplication(localRepository, sqliteRepository, mqttRepository)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	go func() {
		client.Start()
		for {
			select {
			case <-ticker.C:
				if err := application.Process(); err != nil {
					log.Fatal(err)
				}
			case msg := <-client.Message():
				payload := string(msg.Payload()[:])
				elements := strings.Split(msg.Topic(), "/")
				value, _ := strconv.ParseFloat(payload, 64)

				err := application.AddLog(*domain.NewLog(elements, value))
				if err != nil {
					log.Printf(" - error - %s", err)
					time.Sleep(time.Second)
					client.Message() <- msg
				}
			case <-stop:
				log.Println("stopping")
				return
			}
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	stop <- <-done
}
