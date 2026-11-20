package infrastructure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sylvek/domotik/datastore/core/model"
)

// DecodeTasmota transforms a `tele/{device}/SENSOR` MQTT message into a telemetry.
func DecodeTasmota(msg Message) (model.TasmotaTelemetry, error) {
	elements := strings.Split(msg.Topic, "/")
	if len(elements) != 3 {
		return model.TasmotaTelemetry{}, fmt.Errorf("unexpected tasmota topic %s", msg.Topic)
	}

	sensors := map[string]json.RawMessage{}
	if err := json.Unmarshal(msg.Payload, &sensors); err != nil {
		return model.TasmotaTelemetry{}, err
	}

	telemetry := model.TasmotaTelemetry{Device: elements[1], Sensors: map[string]model.TasmotaSensor{}}
	for name, raw := range sensors {
		sensor := model.TasmotaSensor{}
		// fields that are not sensors (Time, TempUnit…) are skipped
		if err := json.Unmarshal(raw, &sensor); err == nil {
			telemetry.Sensors[name] = sensor
		}
	}

	return telemetry, nil
}
