package infrastructure

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sylvek/domotik/datastore/domain/model"
	"github.com/sylvek/domotik/datastore/port"

	_ "github.com/glebarez/go-sqlite"
)

type Aggregation string

const (
	AVG   Aggregation = "AVG(value)"
	MAX   Aggregation = "MAX(value)"
	SUM   Aggregation = "SUM(value)"
	LAST  Aggregation = "(value)"
	DELTA Aggregation = "MAX(value) - MIN(value)"
)

const (
	WATT   string = "watt"
	TEMP   string = "temp"
	LITER  string = "liter"
	RATE   string = "rate"
	INDICE string = "indice"
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

type Parameters struct {
	CurrentIndex int
}

type SqliteClient struct {
	path       string
	parameters Parameters
	instances  map[string]*Instance
}

// Close implements port.LogRepository.
func (d *SqliteClient) Close() {
	for _, v := range d.instances {
		v.db.Close()
	}
}

// Store implements port.LogRepository.
func (d *SqliteClient) Store(l model.Log) error {
	t := time.Now()

	currentDayOfYear := t.YearDay()
	if currentDayOfYear != d.parameters.CurrentIndex {
		d.parameters.CurrentIndex = currentDayOfYear

		log.Printf("daily work is starting (%d)\n", currentDayOfYear)

		d.synchronization()
		d.aggregation()
		d.cleaning()

		log.Println("daily work is finished")
	}

	db := d.instances[l.Topic].db
	_, err := db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (?, ?, ?, ?)",
		// Grafana requires that data is stored in UTC
		t.Unix(),
		l.Name,
		l.Unit,
		l.Value)
	if err != nil {
		return err
	}

	return nil
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
	// We calculate the current offset between localtime and utc0
	// data is stored in utc0. We have to select a slot of data
	// corresponding to 24h in local time.
	diff, _ := strconv.Atoi(time.Now().Format("-07"))
	offset := strconv.Itoa(-1 * diff) // we reverse the value
	for _, v := range d.instances {
		if v.dailyOperations != nil {
			for _, sensorsOperation := range v.dailyOperations {
				sql := fmt.Sprintf("SELECT %s FROM data WHERE name='%s' AND unit='%s'",
					sensorsOperation.aggregate,
					sensorsOperation.from,
					sensorsOperation.unit)
				sql += " AND ts>strftime('%s','now','start of day','" + offset + " hours') AND ts<strftime('%s','now')"
				if sensorsOperation.aggregate == LAST {
					sql += " ORDER BY ts DESC LIMIT 1"
				}
				if err := v.db.QueryRow(sql).Scan(&value); err != nil {
					log.Println("aggregation", "select", sensorsOperation, err)
				} else {
					if _, err := d.instances["history"].db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (strftime('%s','now','start of day'), ?, ?, ?)",
						sensorsOperation.to,
						sensorsOperation.unit,
						value); err != nil {
						log.Println("aggregation", "history", v, err)
					}
				}
			}
		}

	}
}

func (d *SqliteClient) synchronization() {
	p, _ := json.Marshal(d.parameters)
	os.WriteFile(d.path+"/database.json", p, 0644)
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

func NewSqliteDatabase(path string) port.LogRepository {
	parameters := Parameters{}
	data, err := os.ReadFile(path + "/database.json")
	if err == nil {
		json.Unmarshal(data, &parameters)
	}

	instances := make(map[string]*Instance)
	instances["sensors"] = &Instance{
		db:       prepareDatabase(path + "/sensors.db"),
		volatile: true,
		dailyOperations: []Operation{
			{aggregate: AVG, from: "outside", to: "daily_temp_outside", unit: TEMP},
			{aggregate: AVG, from: "living", to: "daily_temp_inside", unit: TEMP},
			{aggregate: DELTA, from: "linky", to: "daily_power_consumption", unit: INDICE},
			{aggregate: LAST, from: "sumPerDay", to: "daily_rate_consumption", unit: RATE}}}
	instances["history"] = &Instance{
		db:       prepareDatabase(path + "/history.db"),
		volatile: false}

	return &SqliteClient{
		path:       path,
		parameters: parameters,
		instances:  instances}
}
