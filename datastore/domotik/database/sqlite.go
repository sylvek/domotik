package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type Aggregation string

const (
	AVG Aggregation = "AVG"
	MAX Aggregation = "MAX"
	SUM Aggregation = "SUM"
)

const (
	WATT  string = "watt"
	TEMP  string = "temp"
	LITER string = "liter"
)

type Operation struct {
	aggregate Aggregation
	from      string
	to        string
	unit      string
}

type Instance struct {
	db              *sql.DB
	dailyOperations []Operation
	volatile        bool
}

type SqliteClient struct {
	databasePath string
	parameters   Parameters
	instances    map[string]*Instance
}

type Parameters struct {
	CurrentIndex int
}

func NewSqliteClient(databasePath string) Database {

	parameters := Parameters{}
	data, err := ioutil.ReadFile(databasePath + "/database.json")
	if err == nil {
		json.Unmarshal(data, &parameters)
	}

	instances := make(map[string]*Instance)
	instances["sensors"] = &Instance{
		db:       prepareDatabase(databasePath + "/sensors.db"),
		volatile: true,
		dailyOperations: []Operation{
			{aggregate: AVG, from: "outside", to: "daily_temp_outside", unit: TEMP},
			{aggregate: AVG, from: "living", to: "daily_temp_inside", unit: TEMP}}}
	instances["measures"] = &Instance{
		db:       prepareDatabase(databasePath + "/measures.db"),
		volatile: true,
		dailyOperations: []Operation{
			{aggregate: MAX, from: "sumPerDay", to: "daily_power_consumption", unit: WATT},
			{aggregate: SUM, from: "waterPerDay", to: "daily_water_consumption", unit: LITER}}}
	instances["history"] = &Instance{
		db:       prepareDatabase(databasePath + "/history.db"),
		volatile: false}

	return &SqliteClient{
		databasePath: databasePath,
		parameters:   parameters,
		instances:    instances}
}

func (d *SqliteClient) Close() {
	for _, v := range d.instances {
		v.db.Close()
	}
}

func (d *SqliteClient) cleaning() {
	for _, v := range d.instances {
		if v.volatile {
			if _, err := v.db.Exec("DELETE FROM data where ts < strftime('%s','now', '-2 day')"); err != nil {
				log.Println("cleaning", err)
			}
		}
	}
}

func (d *SqliteClient) aggregation() {
	var value float64
	for _, kind := range []string{"sensors", "measures"} {
		for _, sensorsOperation := range d.instances[kind].dailyOperations {
			sql := fmt.Sprintf("SELECT %s(value) FROM data WHERE name='%s' AND unit='%s'",
				sensorsOperation.aggregate,
				sensorsOperation.from,
				sensorsOperation.unit)
			if err := d.instances[kind].db.QueryRow(sql + " AND ts>strftime('%s','now','start of day','-1 day') AND ts<strftime('%s','now','start of day')").Scan(&value); err != nil {
				log.Println("aggregation", "select", kind, err)
			} else {
				if _, err := d.instances["history"].db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (strftime('%s','now','start of day','-1 day'), ?, ?, ?)",
					sensorsOperation.to,
					sensorsOperation.unit,
					value); err != nil {
					log.Println("aggregation", "history", kind, err)
				}
			}
		}
	}
}

func (d *SqliteClient) synchronization() {
	p, _ := json.Marshal(d.parameters)
	os.WriteFile(d.databasePath+"/database.json", p, 0644)
}

func (d *SqliteClient) AddSeries(
	topic string,
	name string,
	unit string,
	value float32) error {

	t := time.Now()

	db := d.instances[topic].db
	_, err := db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (?, ?, ?, ?)",
		t.Unix(),
		name,
		unit,
		value)
	if err != nil {
		return err
	}

	currentDayOfYear := t.YearDay()
	if currentDayOfYear != d.parameters.CurrentIndex {
		d.parameters.CurrentIndex = currentDayOfYear

		log.Printf("daily work is starting (%d)\n", currentDayOfYear)

		d.synchronization()
		d.aggregation()
		d.cleaning()

		log.Println("daily work is finished")
	}

	return nil
}

func prepareDatabase(entity string) (db *sql.DB) {
	// https://www.sqlite.org/pragma.html#pragma_busy_timeout
	// https://www.sqlite.org/pragma.html#pragma_journal_mode
	db, err := sql.Open("sqlite", entity+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		panic(err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS data(id INTEGER PRIMARY KEY AUTOINCREMENT, ts INTEGER, name TEXT, unit TEXT, value REAL);")
	db.Exec("CREATE INDEX IF NOT EXISTS indx_data on data (ts);")

	return db
}