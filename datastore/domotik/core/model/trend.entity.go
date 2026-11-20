package model

// Trend compares a value with the value of the same period one year before.
type Trend struct {
	Value float64
	// Previous is nil when the previous period is unknown.
	Previous *float64
}

// MINIMAL_COVERAGE is the share of days the previous period must have,
// compared to the current one, to be comparable.
const MINIMAL_COVERAGE = 0.8

func Sum(points []Point) float64 {
	sum := 0.0
	for _, p := range points {
		sum += p.Value
	}
	return sum
}

func Mean(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	return Sum(points) / float64(len(points))
}

// NewTrend aggregates the daily points of a period and of the same period one year before.
// It returns nil when there is nothing in the current period.
func NewTrend(current []Point, previous []Point, aggregate func([]Point) float64) *Trend {
	if len(current) == 0 {
		return nil
	}
	trend := Trend{Value: aggregate(current)}
	if len(previous) > 0 && float64(len(previous)) >= MINIMAL_COVERAGE*float64(len(current)) {
		value := aggregate(previous)
		trend.Previous = &value
	}
	return &trend
}

// Change is the relative evolution since the previous period (0.1 means +10%).
// It returns nil when it cannot be computed.
func (t Trend) Change() *float64 {
	if t.Previous == nil || *t.Previous == 0 {
		return nil
	}
	change := (t.Value - *t.Previous) / abs(*t.Previous)
	return &change
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// Between returns the points of a period, `from` is included and `to` is excluded.
// Points must be sorted by time.
func Between(points []Point, from int64, to int64) []Point {
	result := []Point{}
	for _, p := range points {
		if p.TS >= from && p.TS < to {
			result = append(result, p)
		}
	}
	return result
}
