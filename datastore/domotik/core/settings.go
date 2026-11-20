package core

// Settings gathers the rules used to build the dashboards.
type Settings struct {
	// Labels gives a human readable name to a sensor (ie. living -> séjour)
	Labels map[string]string
	// OutsideSensor is the name of the sensor placed outside, every other one is considered inside.
	OutsideSensor string
	// WaterEuroByLiter is the price of one liter of water.
	WaterEuroByLiter float64
	// TypicalDays are reference days, sorted by consumption, used to qualify the consumption of a day.
	TypicalDays []TypicalDay
	// A "cold day" is a day where the average outside temperature is below ColdDayTemperature (°C)
	ColdDayTemperature float64
	// A "high consumption" day is a day where more than HighConsumptionWatt (Wh) were consumed.
	HighConsumptionWatt float64
	// A "hot day" is a day where the average outside temperature is above HotDayTemperature (°C)
	HotDayTemperature float64
}

const (
	LEVEL_LOW       = "low"
	LEVEL_MODERATE  = "moderate"
	LEVEL_HIGH      = "high"
	LEVEL_VERY_HIGH = "very high"
)

// TypicalDay is a reference: a kind of day and the electricity (Wh) consumed during it.
type TypicalDay struct {
	Label string
	Watt  int64
	Level string
}

// typicalDay returns the first reference day consuming at least the given consumption.
func (s Settings) typicalDay(wattConsumed int64) TypicalDay {
	for _, day := range s.TypicalDays {
		if wattConsumed <= day.Watt {
			return day
		}
	}
	return TypicalDay{Label: "above top max", Watt: wattConsumed, Level: LEVEL_VERY_HIGH}
}

func NewSettings(labels map[string]string) Settings {
	return Settings{
		Labels:        labels,
		OutsideSensor: "outside",
		// https://github.com/sylvek/wearesober/blob/main/what-ive-learned.md
		TypicalDays: []TypicalDay{
			{Label: "all off (except fridge)", Watt: 3000, Level: LEVEL_LOW},
			{Label: "a winter day, nobody home (frost protection)", Watt: 7000, Level: LEVEL_LOW},
			{Label: "a summer day", Watt: 13000, Level: LEVEL_LOW},
			{Label: "a fall/spring day", Watt: 25000, Level: LEVEL_LOW},
			{Label: "a winter day", Watt: 35000, Level: LEVEL_MODERATE},
			{Label: "top max", Watt: 60000, Level: LEVEL_HIGH},
		},
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
