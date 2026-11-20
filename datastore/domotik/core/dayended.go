package core

import (
	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/port"
)

const DAY_ENDED_MAX_ATTEMPTS = 3

// NotifyDayEnded notifies the electricity consumption of the day that has just ended.
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

	err := a.notifyDayEnded(day, notifier)
	if err != nil {
		// the day is notified again later
		a.mutex.Lock()
		if a.endedDay == 0 && attempts < DAY_ENDED_MAX_ATTEMPTS {
			a.endedDay = day
		}
		a.mutex.Unlock()
	}

	return err
}

func (a *Application) notifyDayEnded(day int64, notifier port.DayEndedNotifier) error {
	consumption, err := a.measureRepository.Daily(DAILY_POWER, day, day+1)
	if err != nil {
		return err
	}
	// nothing was measured during this day
	if len(consumption) == 0 {
		return nil
	}

	return notifier.Notify(dto.DayEnded{
		Timestamp:    day,
		WattConsumed: int64(consumption[len(consumption)-1].Value),
	})
}
