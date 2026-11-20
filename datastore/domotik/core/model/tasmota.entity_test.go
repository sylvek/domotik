package model

import "testing"

func TestTasmotaTelemetryLogs(t *testing.T) {
	temp, humidity, dew := 21.4, 50.1, 10.5
	telemetry := TasmotaTelemetry{Device: "bureau", Sensors: map[string]TasmotaSensor{
		"AM2301": {Temperature: &temp, Humidity: &humidity, DewPoint: &dew},
		"OTHER":  {Temperature: &temp},
	}}

	logs := telemetry.Logs()
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs, got %d", len(logs))
	}
	expected := []struct {
		unit  string
		value float64
	}{{"temp", temp}, {"humidity", humidity}, {"dewpoint", dew}}
	for i, e := range expected {
		if logs[i].GetTopic() != "sensors" || logs[i].GetName() != "bureau" || logs[i].GetUnit() != e.unit || logs[i].GetValue() != e.value {
			t.Errorf("unexpected log %+v", logs[i])
		}
	}

	ds18b20 := TasmotaTelemetry{Device: "cave", Sensors: map[string]TasmotaSensor{"DS18B20": {Temperature: &temp}}}
	if logs := ds18b20.Logs(); len(logs) != 1 || logs[0].GetUnit() != "temp" {
		t.Errorf("unexpected logs %+v", logs)
	}
}
