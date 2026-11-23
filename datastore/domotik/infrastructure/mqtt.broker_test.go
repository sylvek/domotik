package infrastructure

import "testing"

func TestFormatValue(t *testing.T) {
	for value, expected := range map[float64]string{12.7: "12.7", 21: "21", 50.1: "50.1", -3.25: "-3.25", 0: "0"} {
		if got := formatValue(value); got != expected {
			t.Errorf("expected %s, got %s", expected, got)
		}
	}
}
