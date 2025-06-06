package model

type Measure struct {
	lowTariff bool
	index     int64
	timestamp int64
}

func NewMeasure() *Measure {
	return &Measure{lowTariff: false, index: 0, timestamp: 0}
}

func (i *Measure) UpdateFromLog(l Log, timestamp int64) {
	if l.name == "linky" && l.unit == "indice" {
		i.index = int64(l.value)
		i.timestamp = timestamp
	}
	if l.name == "linky" && l.unit == "state" {
		i.lowTariff = l.value == 0.0
	}
}

func (i *Measure) ResetIndex() {
	i.index = 0
}

func (i *Measure) HasNewIndex() bool {
	return i.index > 0
}

func (i *Measure) GetLastIndex() int64 {
	return i.index
}

func (i *Measure) GetLastTimestamp() int64 {
	return i.timestamp
}

func (i *Measure) IsLowTariff() bool {
	return i.lowTariff
}
