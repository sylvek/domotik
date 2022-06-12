package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Aggregation string

const (
	AVG Aggregation = "AVG"
	MAX Aggregation = "MAX"
)

const (
	WATT string = "watt"
	TEMP string = "temp"
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

type Database struct {
	databasePath string
	parameters   Parameters
	instances    map[string]*Instance
}

type Parameters struct {
	CurrentIndex int
}

func NewDatabase(databasePath string) Database {

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
			{aggregate: AVG, from: "esp12e", to: "daily_temp_outside", unit: TEMP},
			{aggregate: AVG, from: "esp8266", to: "daily_temp_inside", unit: TEMP}}}
	instances["measures"] = &Instance{
		db:       prepareDatabase(databasePath + "/measures.db"),
		volatile: true,
		dailyOperations: []Operation{
			{aggregate: MAX, from: "sumPerDay", to: "daily_power_consumption", unit: WATT}}}
	instances["history"] = &Instance{
		db:       prepareDatabase(databasePath + "/history.db"),
		volatile: false}

	return Database{databasePath: databasePath, parameters: parameters, instances: instances}
}

func (d *Database) Close() {
	for _, v := range d.instances {
		v.db.Close()
	}
}

func (d *Database) cleaning() {
	for _, v := range d.instances {
		if v.volatile {
			v.db.Exec("DELETE FROM data where ts < strftime('%s','now', '-2 day')")
		}
	}
}

func (d *Database) aggregation() {
	var value float64
	for _, kind := range []string{"sensors", "measures"} {
		for _, sensorsOperation := range d.instances[kind].dailyOperations {
			sql := fmt.Sprintf("SELECT %s(value) FROM data WHERE name='%s'", sensorsOperation.aggregate, sensorsOperation.from)
			err := d.instances[kind].db.QueryRow(sql + " AND ts>strftime('%s','now','start of day','-1 day') AND ts<strftime('%s','now','start of day')").Scan(&value)

			if err == nil {
				d.instances["history"].db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (strftime('%s','now','start of day','-1 day'), ?, ?, ?)",
					sensorsOperation.to,
					sensorsOperation.unit,
					value)
			}
		}
	}
}

func (d *Database) synchronization() {
	p, _ := json.Marshal(d.parameters)
	os.WriteFile(d.databasePath+"/database.json", p, 0644)
}

func (d *Database) AddSeries(
	kind string,
	name string,
	unit string,
	value float32) error {

	t := time.Now()

	log.Printf("%s: %s -> %f (%s)\n", kind, name, value, unit)

	db := d.instances[kind].db
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

		log.Println("daily work is starting")

		d.synchronization()
		d.aggregation()
		d.cleaning()

		log.Println("daily work is finished")
	}

	return nil
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
