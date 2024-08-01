package infrastructure

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client   mqtt.Client
	messages chan mqtt.Message
}

func NewMQTTClient(broker string, topics []string) *MqttClient {

	messages := make(chan mqtt.Message, 10)

	opts := mqtt.NewClientOptions().AddBroker(broker).SetOrderMatters(false)
	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("status: connected")
		for _, topic := range topics {
			if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
			log.Printf("topic [%s] subscribed", topic)
		}
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) { log.Println("status: connection lost") }
	opts.OnReconnecting = func(c mqtt.Client, co *mqtt.ClientOptions) { log.Println("status: reconnecting") }
	opts.DefaultPublishHandler = func(client mqtt.Client, msg mqtt.Message) {
		messages <- msg
	}

	return &MqttClient{client: mqtt.NewClient(opts), messages: messages}
}

func (mqtt *MqttClient) Message() chan mqtt.Message {
	return mqtt.messages
}

func (m *MqttClient) Publish(topic string, name string, value float64, unit string) {
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

func (mqtt *MqttClient) Stop() error {
	token := mqtt.client.Unsubscribe()
	if token.Error() != nil {
		return token.Error()
	}
	_ = token.Wait()
	mqtt.client.Disconnect(250)
	return nil
}

func (mqtt *MqttClient) Start() error {
	token := mqtt.client.Connect()
	if token.Error() != nil {
		return token.Error()
	}
	_ = token.Wait()
	return nil
}
