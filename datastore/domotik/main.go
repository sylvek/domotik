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

type Database struct {
	Db       *sql.DB
	Volatile bool
}

func main() {
	mqttHost := os.Getenv("MQTT_HOST")
	databasePath := os.Getenv("DB_PATH")
	log.Printf("starting - MQTT_HOST:%s DB_PATH:%s", mqttHost, databasePath)

	// prepare database
	databases := make(map[string]Database)
	databases["sensors"] = Database{Db: prepareDatabase(databasePath + "/sensors.db"), Volatile: true}
	databases["measures"] = Database{Db: prepareDatabase(databasePath + "/measures.db"), Volatile: false}

	// connect
	client, logs := connect(mqttHost, []string{"sensors", "measures"})

	go treatLogs(logs, databases)

	defer func() {
		databases["sensors"].Db.Close()
		databases["measures"].Db.Close()
		client.Disconnect(250)
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Println("ciao.")
}

func prepareDatabase(entity string) (db *sql.DB) {
	db, err := sql.Open("sqlite", entity)
	if err != nil {
		panic(err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS data(id INTEGER PRIMARY KEY AUTOINCREMENT, ts INTEGER, name TEXT, unit TEXT, value REAL);")
	db.Exec("CREATE INDEX IF NOT EXISTS indx_data on data (ts);")

	return db
}

func connect(mqttBroker string, kinds []string) (client mqtt.Client, logs chan Log) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBroker, 1883))
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload()[:])
		elements := strings.Split(msg.Topic(), "/")
		value, _ := strconv.ParseFloat(payload, 32)
		logs <- Log{Kind: elements[0], Name: elements[1], Unit: elements[2], Value: float32(value)}
	})
	logs = make(chan Log, 10)
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, kind := range kinds {
		client.Subscribe(kind+"/#", 1, nil)
	}
	return client, logs
}

func treatLogs(logs chan Log, databases map[string]Database) {
	previousDay := -1
	for {
		t := time.Now()
		currentDay := t.YearDay()

		item := <-logs
		log.Printf("%s: %s -> %f (%s)\n", item.Kind, item.Name, item.Value, item.Unit)

		db := databases[item.Kind].Db
		volatile := databases[item.Kind].Volatile
		if currentDay != previousDay && volatile {
			previousDay = currentDay
			epoch := t.Unix()

			// we clean old data (meaning older than 2days)
			db.Exec(fmt.Sprintf("DELETE FROM data where ts < %d", epoch-172800 /* 2days */))
		}
		_, err := db.Exec(fmt.Sprintf("INSERT INTO data (ts, name, unit, value) VALUES (%d, '%s', '%s', %f);",
			t.Unix(),
			item.Name,
			item.Unit,
			item.Value))
		if err != nil {
			log.Printf(" - error - %s", err)
			time.Sleep(time.Second)
			logs <- item
		}
	}
}
