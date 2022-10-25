package broker

import (
	"fmt"
	"strconv"
	"strings"

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
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload()[:])
		elements := strings.Split(msg.Topic(), "/")
		value, _ := strconv.ParseFloat(payload, 32)
		logs <- Log{Topic: elements[0], Name: elements[1], Unit: elements[2], Value: float32(value)}
	})
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, topic := range c.topics {
		if token := c.client.Subscribe(topic+"/+/+", 1, nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
}

// Disconnect implements BrokerClient
func (c *MqttClient) Disconnect() {
	c.client.Disconnect(250)
}

func NewMQTTBrokerClient(host string, port int) BrokerClient {
	return &MqttClient{mqttBroker: fmt.Sprintf("tcp://%s:%d", host, port), topics: []string{"sensors", "measures"}}
}
