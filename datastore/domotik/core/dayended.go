package core

import (
	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/port"
)

const DAY_ENDED_MAX_ATTEMPTS = 3

// NotifyDayEnded notifies the measures of the day that has just ended.
// It does nothing while no day is waiting for its notification.
func (a *Application) NotifyDayEnded(notifier port.DayEndedNotifier) error {
	a.mutex.Lock()
	day := a.endedDay
	if day == 0 {
		a.mutex.Unlock()
		return nil
	}
	a.endedDay = 0
	a.endedDayAttempts++
	attempts := a.endedDayAttempts
	a.mutex.Unlock()

	for _, name := range a.settings.DayEndedMeasures {
		value, err := a.measureRepository.DailyValue(name, day)
		if err == nil && value != nil {
			err = notifier.Notify(dto.DayEnded{Timestamp: value.TS, Name: value.Name, Unit: value.Unit, Value: value.Value})
		}
		if err != nil {
			// the whole day is notified again later
			a.mutex.Lock()
			if a.endedDay == 0 && attempts < DAY_ENDED_MAX_ATTEMPTS {
				a.endedDay = day
			}
			a.mutex.Unlock()
			return err
		}
	}

	return nil
}
