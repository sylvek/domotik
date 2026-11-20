package infrastructure

import (
	"fmt"
	"io"
	"log"
	"log/slog"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

type Message struct {
	Topic   string
	Payload []byte
}

type outgoing struct {
	topic   string
	payload []byte
}

// MqttBroker is an embedded MQTT broker (TCP listener).
// The application talks to it in-process: no external broker is required.
type MqttBroker struct {
	server   *mqtt.Server
	topics   map[string]chan Message
	outgoing chan outgoing
	done     chan struct{}
}

func NewMQTTBroker(address string) (*MqttBroker, error) {
	server := mqtt.New(&mqtt.Options{
		InlineClient: true,
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	if err := server.AddHook(new(auth.AllowHook), nil); err != nil {
		return nil, err
	}
	if err := server.AddListener(listeners.NewTCP(listeners.Config{ID: "tcp", Address: address})); err != nil {
		return nil, err
	}

	return &MqttBroker{
		server:   server,
		topics:   make(map[string]chan Message),
		outgoing: make(chan outgoing, 100),
		done:     make(chan struct{}),
	}, nil
}

// Messages returns the channel receiving every message published on the given topic filter.
// It must be called before Start.
func (m *MqttBroker) Messages(filter string) chan Message {
	if _, ok := m.topics[filter]; !ok {
		m.topics[filter] = make(chan Message, 10)
	}
	return m.topics[filter]
}

// Publish queues a retained message. The delivery is done by a dedicated goroutine
// so that the caller is never blocked by its own subscriptions.
func (m *MqttBroker) Publish(topic string, name string, value float64, unit string) {
	select {
	case m.outgoing <- outgoing{
		topic:   fmt.Sprintf("%s/%s/%s", topic, name, unit),
		payload: []byte(fmt.Sprintf("%.3f", value))}:
	default:
		log.Println("unable to publish, queue is full", topic, name, unit)
	}
}

func (m *MqttBroker) Start() error {
	id := 1
	for filter, messages := range m.topics {
		err := m.server.Subscribe(filter, id, func(_ *mqtt.Client, _ packets.Subscription, pk packets.Packet) {
			payload := make([]byte, len(pk.Payload))
			copy(payload, pk.Payload)
			select {
			case messages <- Message{Topic: pk.TopicName, Payload: payload}:
			case <-m.done:
			}
		})
		if err != nil {
			return err
		}
		log.Printf("topic [%s] subscribed", filter)
		id++
	}

	go func() {
		for {
			select {
			case o := <-m.outgoing:
				if err := m.server.Publish(o.topic, o.payload, true, 0); err != nil {
					log.Println("unable to publish", o.topic, err)
				}
			case <-m.done:
				return
			}
		}
	}()

	return m.server.Serve()
}

func (m *MqttBroker) Stop() error {
	close(m.done)
	return m.server.Close()
}
