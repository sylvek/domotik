package core

import (
	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/port"
)

// notifyDayEnded notifies, once, the electricity consumption of the day that has just ended.
// It does nothing while no day is waiting for its notification.
func (a *Application) notifyDayEnded(notifier port.DayEndedNotifier) error {
	day := a.endedDay
	if day == 0 {
		return nil
	}
	a.endedDay = 0

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
