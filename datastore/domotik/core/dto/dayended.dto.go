package dto

// DayEnded describes a day that has just ended.
type DayEnded struct {
	// Timestamp is midnight UTC of the ended day, in seconds.
	Timestamp int64
	// WattConsumed is the electricity consumed during the day, in Wh.
	WattConsumed int64
}
