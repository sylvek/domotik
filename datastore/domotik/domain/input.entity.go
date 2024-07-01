package domain

type Input struct {
	lowTariff bool
	indice    int64
}

func NewInput() *Input {
	return &Input{lowTariff: false, indice: 0}
}

func (i *Input) UpdateFromLog(l Log) {
	if l.name == "linky" && l.unit == "indice" {
		i.indice = int64(l.value)
	}
	if l.name == "linky" && l.unit == "state" {
		i.lowTariff = l.value == 0.0
	}
}

func (i *Input) ResetIndice() {
	i.indice = 0
}

func (i *Input) HasIndice() bool {
	return i.indice > 0
}

func (i *Input) GetIndice() int64 {
	return i.indice
}

func (i *Input) IsLowTariff() bool {
	return i.lowTariff
}
