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
	kind  string
	name  string
	unit  string
	value float32
}

type Aggregation string

const (
	AVG Aggregation = "AVG"
	SUM Aggregation = "SUM"
)

const (
	WATT string = "watt"
	TEMP string = "temp"
)

type Operation struct {
	aggregat Aggregation
	from     string
	to       string
	unit     string
}

type Database struct {
	db              *sql.DB
	dailyOperations []Operation
}

func main() {
	mqttHost := os.Getenv("MQTT_HOST")
	databasePath := os.Getenv("DB_PATH")
	log.Printf("starting - MQTT_HOST:%s DB_PATH:%s", mqttHost, databasePath)

	// prepare database
	databases := make(map[string]Database)
	databases["sensors"] = Database{db: prepareDatabase(databasePath + "/sensors.db"), dailyOperations: []Operation{
		{aggregat: AVG, from: "esp12e", to: "daily_temp_outside", unit: TEMP},
		{aggregat: AVG, from: "esp8266", to: "daily_temp_inside", unit: TEMP}}}
	databases["measures"] = Database{db: prepareDatabase(databasePath + "/measures.db"), dailyOperations: []Operation{
		{aggregat: SUM, from: "meanPerHour", to: "daily_power_consumption", unit: WATT}}}
	databases["history"] = Database{db: prepareDatabase(databasePath + "/history.db")}

	// connect
	client, logs := connect(mqttHost, []string{"sensors", "measures"})

	go treatLogs(logs, databases)

	defer func() {
		databases["sensors"].db.Close()
		databases["measures"].db.Close()
		databases["history"].db.Close()
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
		logs <- Log{kind: elements[0], name: elements[1], unit: elements[2], value: float32(value)}
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

		if currentDay != previousDay {
			previousDay = currentDay
			epoch := t.Unix()

			// we clean old data (meaning older than 2days)
			databases["sensors"].db.Exec("DELETE FROM data where ts < ?", epoch-172800 /* 2days */)
			databases["measures"].db.Exec("DELETE FROM data where ts < ?", epoch-172800 /* 2days */)

			// aggregate data
			var value float64
			for _, kind := range []string{"sensors", "measures"} {
				for _, sensorsOperation := range databases[kind].dailyOperations {
					databases["sensors"].db.QueryRow(fmt.Sprintf("SELECT %s(value) FROM data WHERE name='%s' AND ts>%d", sensorsOperation.aggregat, sensorsOperation.from, 0)).Scan(&value)
					databases["history"].db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (?, ?, ?, ?)",
						0,
						sensorsOperation.to,
						sensorsOperation.unit,
						value)
				}
			}
		}

		item := <-logs
		log.Printf("%s: %s -> %f (%s)\n", item.kind, item.name, item.value, item.unit)

		db := databases[item.kind].db
		_, err := db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (?, ?, ?, ?)",
			t.Unix(),
			item.name,
			item.unit,
			item.value)
		if err != nil {
			log.Printf(" - error - %s", err)
			time.Sleep(time.Second)
			logs <- item
		}
	}
}
