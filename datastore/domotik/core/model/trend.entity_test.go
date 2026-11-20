package model

import (
	"testing"
	"time"
)

func TestSeasonOf(t *testing.T) {
	expected := map[time.Month]Season{
		time.January: WINTER, time.March: SPRING, time.May: SPRING, time.June: SUMMER,
		time.August: SUMMER, time.September: AUTUMN, time.November: AUTUMN, time.December: WINTER,
	}
	for month, season := range expected {
		if got := SeasonOf(ts(2025, month, 15)); got != season {
			t.Errorf("%s: expected %s, got %s", month, season, got)
		}
	}
}

func TestTrend(t *testing.T) {
	current := []Point{{TS: 1, Value: 10}, {TS: 2, Value: 12}}
	previous := []Point{{TS: 1, Value: 10}, {TS: 2, Value: 10}}

	trend := NewTrend(current, previous, Sum)
	if trend.Value != 22 || *trend.Previous != 20 || *trend.Change() != 0.1 {
		t.Errorf("unexpected trend %+v", trend)
	}

	// the previous period is not comparable when too many days are missing
	trend = NewTrend(append(current, current...), previous[:1], Sum)
	if trend.Previous != nil || trend.Change() != nil {
		t.Errorf("unexpected trend %+v", trend)
	}

	if NewTrend(nil, previous, Sum) != nil {
		t.Error("no trend was expected without current value")
	}

	// a temperature can be negative
	trend = NewTrend([]Point{{Value: -2}}, []Point{{Value: -4}}, Mean)
	if *trend.Change() != 0.5 {
		t.Errorf("unexpected change %v", *trend.Change())
	}
}
