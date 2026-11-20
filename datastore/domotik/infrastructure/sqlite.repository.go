package infrastructure

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sylvek/domotik/datastore/core/model"
	"github.com/sylvek/domotik/datastore/core/port"

	_ "github.com/glebarez/go-sqlite"
)

var _ port.LogRepository = (*SqliteDatabase)(nil)
var _ port.MeasureRepository = (*SqliteDatabase)(nil)

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
	EURO   string = "euro"
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

type SqliteDatabase struct {
	path       string
	parameters Parameters
	instances  map[string]*Instance
}

// Close implements port.LogRepository.
func (d *SqliteDatabase) Close() {
	for _, v := range d.instances {
		v.db.Close()
	}
}

// CloseDay implements port.LogRepository.
func (d *SqliteDatabase) CloseDay() bool {
	currentDayOfYear := time.Now().YearDay()
	if currentDayOfYear == d.parameters.CurrentIndex {
		return false
	}
	d.parameters.CurrentIndex = currentDayOfYear

	log.Printf("daily work is starting (%d)\n", currentDayOfYear)

	d.synchronization()
	d.aggregation()
	d.cleaning()

	log.Println("daily work is finished")

	return true
}

// Store implements port.LogRepository.
func (d *SqliteDatabase) Store(l model.Log) error {
	t := time.Now()

	db := d.instances[l.GetTopic()].db
	_, err := db.Exec("INSERT INTO data (ts, name, unit, value) VALUES (?, ?, ?, ?)",
		// data is stored in UTC
		t.Unix(),
		l.GetName(),
		l.GetUnit(),
		l.GetValue())
	if err != nil {
		return err
	}

	return nil
}

func (d *SqliteDatabase) cleaning() {
	for _, v := range d.instances {
		if v.volatile {
			if _, err := v.db.Exec("DELETE FROM data where ts < strftime('%s','now', '-2 day')"); err != nil {
				log.Println("cleaning", err)
			}
		}
	}
}

func (d *SqliteDatabase) aggregation() {
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

func (d *SqliteDatabase) synchronization() {
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
	db.Exec("CREATE INDEX IF NOT EXISTS indx_data_series on data (name, unit, ts);")

	return db
}

// Last implements port.MeasureRepository.
func (d *SqliteDatabase) Last(name string, unit string, from int64, to int64) (*model.Point, error) {
	point := model.Point{}
	err := d.instances["sensors"].db.QueryRow(
		"SELECT ts, value FROM data WHERE name=? AND unit=? AND ts>=? AND ts<? ORDER BY ts DESC LIMIT 1",
		name, unit, from, to).Scan(&point.TS, &point.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &point, nil
}

// Series implements port.MeasureRepository.
func (d *SqliteDatabase) Series(name string, unit string, from int64, to int64) (map[string][]model.Point, error) {
	rows, err := d.instances["sensors"].db.Query(
		"SELECT name, ts, value FROM data WHERE (name=? OR ?='') AND unit=? AND ts>=? AND ts<? ORDER BY ts",
		name, name, unit, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	series := map[string][]model.Point{}
	for rows.Next() {
		var sensor string
		point := model.Point{}
		if err := rows.Scan(&sensor, &point.TS, &point.Value); err != nil {
			return nil, err
		}
		series[sensor] = append(series[sensor], point)
	}
	return series, rows.Err()
}

// Daily implements port.MeasureRepository.
func (d *SqliteDatabase) Daily(name string, from int64, to int64) ([]model.Point, error) {
	return d.daily("SELECT ts, value FROM data WHERE name=? AND ts>=? AND ts<? AND value IS NOT NULL ORDER BY ts",
		name, from, to)
}

// DailyValue implements port.MeasureRepository.
func (d *SqliteDatabase) DailyValue(name string, day int64) (*model.DailyValue, error) {
	value := model.DailyValue{Name: name}
	err := d.instances["history"].db.QueryRow(
		"SELECT ts, unit, value FROM data WHERE name=? AND ts=? AND value IS NOT NULL ORDER BY id DESC LIMIT 1",
		name, day).Scan(&value.TS, &value.Unit, &value.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// DailySum implements port.MeasureRepository.
func (d *SqliteDatabase) DailySum(unit string, from int64, to int64) ([]model.Point, error) {
	return d.daily("SELECT ts - ts % 86400 AS day, SUM(value) FROM data WHERE unit=? AND ts>=? AND ts<? AND value IS NOT NULL GROUP BY day ORDER BY day",
		unit, from, to)
}

func (d *SqliteDatabase) daily(query string, args ...any) ([]model.Point, error) {
	rows, err := d.instances["history"].db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := []model.Point{}
	for rows.Next() {
		point := model.Point{}
		if err := rows.Scan(&point.TS, &point.Value); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, rows.Err()
}

func NewSqliteDatabase(path string) *SqliteDatabase {
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
			{aggregate: LAST, from: "sumPerDay", to: "daily_rate_consumption", unit: RATE},
			{aggregate: LAST, from: "sumPerDay", to: "daily_cost_consumption", unit: EURO}}}
	instances["history"] = &Instance{
		db:       prepareDatabase(path + "/history.db"),
		volatile: false}

	return &SqliteDatabase{
		path:       path,
		parameters: parameters,
		instances:  instances}
}
