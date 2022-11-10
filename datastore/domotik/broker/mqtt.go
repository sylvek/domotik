package broker

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client     mqtt.Client
	mqttBroker string
	topics     []string
}

// ConnectAndListen implements BrokerClient
func (c *MqttClient) ConnectAndListen(logs chan Log) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.mqttBroker)
	opts.SetOrderMatters(false)
	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("status: connected")
		for _, topic := range c.topics {
			if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
			log.Printf("topic [%s] subscribed", topic)
		}
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) { log.Println("status: connection lost") }
	opts.OnReconnecting = func(c mqtt.Client, co *mqtt.ClientOptions) { log.Println("status: reconnecting") }
	opts.DefaultPublishHandler = func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload()[:])
		elements := strings.Split(msg.Topic(), "/")
		value, _ := strconv.ParseFloat(payload, 32)
		logs <- Log{Topic: elements[0], Name: elements[1], Unit: elements[2], Value: float32(value)}
	}

	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// Disconnect implements BrokerClient
func (c *MqttClient) Disconnect() {
	c.client.Disconnect(250)
}

func NewMQTTBrokerClient(host string, port int) BrokerClient {
	return &MqttClient{mqttBroker: fmt.Sprintf("tcp://%s:%d", host, port), topics: []string{"sensors/+/+"}}
}
