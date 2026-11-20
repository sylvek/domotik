package dto

// DayEnded is the value of a daily measure for a day that has just ended.
type DayEnded struct {
	// Timestamp is midnight UTC of the ended day, in seconds.
	Timestamp int64
	Name      string
	Unit      string
	Value     float64
}
