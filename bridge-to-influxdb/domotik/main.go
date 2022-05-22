package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func main() {

	log.Printf("starting : %s", os.Args)
	client := connect(os.Args[1], os.Args[2])

	sub(client)
	defer client.Disconnect(250)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Println("ciao.")
}

func connect(mqttBroker string, influxDbHost string) (client mqtt.Client) {
	influx := influxdb2.NewClient(fmt.Sprintf("http://%s:%d", influxDbHost, 8086), "")
	writeAPI := influx.WriteAPIBlocking("", "domotik")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBroker, 1883))
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		payload := string(msg.Payload()[:])
		elements := strings.Split(msg.Topic(), "/")
		value, _ := strconv.ParseFloat(payload, 32)
		p := influxdb2.NewPoint(elements[0],
			map[string]string{"sensor": elements[1]},
			map[string]interface{}{"value": value, "type": elements[0], "name": elements[1], "unit": elements[2]},
			time.Now())
		_ = writeAPI.WritePoint(context.Background(), p)
	})
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return
}

func sub(client mqtt.Client) {
	token1 := client.Subscribe("sensors/#", 1, nil)
	token1.Wait()
	token2 := client.Subscribe("measures/#", 1, nil)
	token2.Wait()
}
