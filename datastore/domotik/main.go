package main

import (
	"database/sql"
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
)

type Log struct {
	Kind  string
	Name  string
	Unit  string
	Value float32
}

func main() {

	log.Printf("starting : %s", os.Args)

	// open+create database
	db, err := sql.Open("sqlite", os.Args[2])
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS measures(id INTEGER PRIMARY KEY AUTOINCREMENT, ts INTEGER, name TEXT, unit TEXT, value REAL);")
	db.Exec("CREATE TABLE IF NOT EXISTS sensors(id INTEGER PRIMARY KEY AUTOINCREMENT, ts INTEGER, name TEXT, unit TEXT, value REAL);")
	db.Exec("CREATE INDEX IF NOT EXISTS indx_measures on measures (ts);")
	db.Exec("CREATE INDEX IF NOT EXISTS indx_sensors on sensors (ts);")

	// connect
	logs := make(chan Log, 10)
	client := connect(os.Args[1], logs)

	go func() {
		for {
			item := <-logs
			_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (ts, name, unit, value) VALUES (%d, '%s', '%s', %f);",
				item.Kind,
				time.Now().Unix(),
				item.Name,
				item.Unit,
				item.Value))
			if err != nil {
				log.Printf(" - error - %s", err)
				time.Sleep(time.Second)
				logs <- item
			}
		}
	}()

	sub(client)
	defer client.Disconnect(250)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Println("ciao.")
}

func connect(mqttBroker string, logs chan Log) (client mqtt.Client) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBroker, 1883))
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		payload := string(msg.Payload()[:])
		elements := strings.Split(msg.Topic(), "/")
		value, _ := strconv.ParseFloat(payload, 32)
		logs <- Log{Kind: elements[0], Name: elements[1], Unit: elements[2], Value: float32(value)}
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
