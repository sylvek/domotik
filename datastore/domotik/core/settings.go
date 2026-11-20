package core

// Settings gathers the rules used to build the dashboards.
type Settings struct {
	// Labels gives a human readable name to a sensor (ie. living -> séjour)
	Labels map[string]string
	// OutsideSensor is the name of the sensor placed outside, every other one is considered inside.
	OutsideSensor string
	// WaterEuroByLiter is the price of one liter of water.
	WaterEuroByLiter float64
	// DayEndedMeasures are the daily measures notified when a day ends.
	DayEndedMeasures []string
	// A "cold day" is a day where the average outside temperature is below ColdDayTemperature (°C)
	ColdDayTemperature float64
	// A "high consumption" day is a day where more than HighConsumptionWatt (Wh) were consumed.
	HighConsumptionWatt float64
	// A "hot day" is a day where the average outside temperature is above HotDayTemperature (°C)
	HotDayTemperature float64
}

func NewSettings(labels map[string]string, dayEndedMeasures []string) Settings {
	return Settings{
		Labels:              labels,
		OutsideSensor:       "outside",
		DayEndedMeasures:    dayEndedMeasures,
		WaterEuroByLiter:    0.004,
		ColdDayTemperature:  5,
		HighConsumptionWatt: 40000,
		HotDayTemperature:   25,
	}
}

func (s Settings) label(name string) string {
	if label, ok := s.Labels[name]; ok {
		return label
	}
	return name
}
