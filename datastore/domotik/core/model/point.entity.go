package model

import (
	"fmt"
	"time"
)

// Point is a value measured at a given time (unix timestamp, in seconds).
type Point struct {
	TS    int64
	Value float64
}

// Downsample averages points by slots of `seconds` seconds.
// Points must be sorted by time.
func Downsample(points []Point, seconds int64) []Point {
	result := []Point{}
	var slot, n int64
	var sum float64
	flush := func() {
		if n > 0 {
			result = append(result, Point{TS: slot, Value: sum / float64(n)})
		}
	}
	for _, p := range points {
		s := p.TS - p.TS%seconds
		if s != slot {
			flush()
			slot, n, sum = s, 0, 0
		}
		sum += p.Value
		n++
	}
	flush()
	return result
}

// Bucket is an aggregation of points sharing the same label (a week, a month, a year…)
type Bucket struct {
	// TS is the time of the first point of the bucket.
	TS    int64
	Label string
	Value float64
}

// Daily points are stamped at midnight UTC of the day they describe.
func day(p Point) time.Time {
	return time.Unix(p.TS, 0).UTC()
}

func ByWeek(p Point) (string, string) {
	year, week := day(p).ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week), fmt.Sprintf("week %02d", week)
}

func ByMonth(p Point) (string, string) {
	month := day(p).Format("2006-01")
	return month, month
}

func ByYear(p Point) (string, string) {
	year := day(p).Format("2006")
	return year, "year " + year
}

// GroupSum sums points by group. `by` returns the key and the label of the group of a point.
// Points must be sorted by time.
func GroupSum(points []Point, by func(Point) (string, string)) []Bucket {
	buckets := []Bucket{}
	current := ""
	for _, p := range points {
		key, label := by(p)
		if key != current || len(buckets) == 0 {
			buckets = append(buckets, Bucket{TS: p.TS, Label: label})
			current = key
		}
		buckets[len(buckets)-1].Value += p.Value
	}
	return buckets
}

// GroupCount counts the matching points by group. Groups without any matching point are kept.
// Points must be sorted by time.
func GroupCount(points []Point, by func(Point) (string, string), match func(Point) bool) []Bucket {
	matches := make([]Point, len(points))
	for i, p := range points {
		matches[i] = Point{TS: p.TS}
		if match(p) {
			matches[i].Value = 1
		}
	}
	return GroupSum(matches, by)
}
