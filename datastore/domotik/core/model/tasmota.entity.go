package model

type TasmotaSensor struct {
	Temperature *float64
	Humidity    *float64
	DewPoint    *float64
}

// TasmotaTelemetry is the content of a `tele/{device}/SENSOR` message.
type TasmotaTelemetry struct {
	Device  string
	Sensors map[string]TasmotaSensor
}

// Logs converts the telemetry into `sensors/{device}/{unit}` logs.
func (t TasmotaTelemetry) Logs() []Log {
	logs := []Log{}
	add := func(unit string, value *float64) {
		if value != nil {
			logs = append(logs, *NewLog([]string{"sensors", t.Device, unit}, *value))
		}
	}

	if sensor, ok := t.Sensors["DS18B20"]; ok {
		add("temp", sensor.Temperature)
	}
	if sensor, ok := t.Sensors["AM2301"]; ok {
		add("temp", sensor.Temperature)
		add("humidity", sensor.Humidity)
		add("dewpoint", sensor.DewPoint)
	}

	return logs
}
