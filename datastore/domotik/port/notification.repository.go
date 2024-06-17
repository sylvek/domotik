package port

import "github.com/sylvek/domotik/datastore/domain"

type NotificationRepository interface {
	Repository
	Notify(domain.Output) error
}
