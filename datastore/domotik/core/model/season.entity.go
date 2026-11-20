package model

import "time"

// Season is a meteorological season of the northern hemisphere.
type Season string

const (
	WINTER Season = "winter"
	SPRING Season = "spring"
	SUMMER Season = "summer"
	AUTUMN Season = "autumn"
)

// SeasonOf returns the season of a day (unix timestamp in seconds, UTC):
// winter starts in december, spring in march, summer in june and autumn in september.
func SeasonOf(ts int64) Season {
	switch time.Unix(ts, 0).UTC().Month() {
	case time.December, time.January, time.February:
		return WINTER
	case time.March, time.April, time.May:
		return SPRING
	case time.June, time.July, time.August:
		return SUMMER
	}
	return AUTUMN
}
