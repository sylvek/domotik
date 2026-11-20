package model

import (
	"testing"
	"time"
)

func ts(year int, month time.Month, day int) int64 {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix()
}

func TestDownsample(t *testing.T) {
	points := Downsample([]Point{{TS: 0, Value: 1}, {TS: 60, Value: 3}, {TS: 300, Value: 10}}, 300)
	if len(points) != 2 || points[0] != (Point{TS: 0, Value: 2}) || points[1] != (Point{TS: 300, Value: 10}) {
		t.Errorf("unexpected points %+v", points)
	}
}

func TestGroupSum(t *testing.T) {
	points := []Point{
		{TS: ts(2025, 12, 28), Value: 1}, // sunday, week 52
		{TS: ts(2025, 12, 29), Value: 2}, // monday, week 1 of 2026
		{TS: ts(2026, 1, 2), Value: 3},
	}

	weeks := GroupSum(points, ByWeek)
	if len(weeks) != 2 || weeks[0] != (Bucket{TS: points[0].TS, Label: "week 52", Value: 1}) || weeks[1] != (Bucket{TS: points[1].TS, Label: "week 01", Value: 5}) {
		t.Errorf("unexpected weeks %+v", weeks)
	}
	months := GroupSum(points, ByMonth)
	if len(months) != 2 || months[0] != (Bucket{TS: points[0].TS, Label: "2025-12", Value: 3}) || months[1] != (Bucket{TS: points[2].TS, Label: "2026-01", Value: 3}) {
		t.Errorf("unexpected months %+v", months)
	}
}

func TestGroupCount(t *testing.T) {
	points := []Point{{TS: ts(2024, 7, 1), Value: 20}, {TS: ts(2025, 7, 1), Value: 30}, {TS: ts(2025, 7, 2), Value: 31}}
	years := GroupCount(points, ByYear, func(p Point) bool { return p.Value > 25 })
	if len(years) != 2 || years[0] != (Bucket{TS: points[0].TS, Label: "year 2024", Value: 0}) || years[1] != (Bucket{TS: points[1].TS, Label: "year 2025", Value: 2}) {
		t.Errorf("unexpected years %+v", years)
	}
}
