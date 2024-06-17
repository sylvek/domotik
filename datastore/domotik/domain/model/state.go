package model

type State struct {
	LastIndice      int64
	LastIndiceTS    int64
	CurrentDay      int
	CurrentHour     int
	DailySumHigh    int64
	DailySumLow     int64
	HourlySum       int64
	HourlyNbIndices int64
}
