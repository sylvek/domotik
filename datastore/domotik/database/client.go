package database

type Database interface {
	Close()
	AddSeries(topic string, name string, unit string, value float64) error
}
