package port

import "github.com/sylvek/domotik/datastore/domain/model"

type NotificationRepository interface {
	Repository
	Notify(model.Output) error
}
