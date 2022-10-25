package database

type Database interface {
	Close()
	AddSeries(topic string, name string, unit string, value float32) error
}
