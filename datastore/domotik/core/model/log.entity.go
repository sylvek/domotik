package model

type Log struct {
	topic string
	name  string
	unit  string
	value float64
}

func NewLog(elements []string, value float64) *Log {
	return &Log{
		topic: elements[0],
		name:  elements[1],
		unit:  elements[2],
		value: value,
	}
}

func (l *Log) GetTopic() string {
	return l.topic
}

func (l *Log) GetName() string {
	return l.name
}

func (l *Log) GetUnit() string {
	return l.unit
}

func (l *Log) GetValue() float64 {
	return l.value
}
