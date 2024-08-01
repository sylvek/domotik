package infrastructure

import (
	"encoding/json"
	"os"

	"github.com/sylvek/domotik/datastore/core/model"
	"github.com/sylvek/domotik/datastore/core/port"
)

type LocalStore struct {
	path string
}

const FILE = "/state.json"

// Retrieve implements port.StateRepository.
func (l *LocalStore) Retrieve() (model.State, error) {
	state := model.NewState()
	data, err := os.ReadFile(l.path + FILE)
	if err == nil {
		json.Unmarshal(data, &state)
	}

	// if the file does not exist we return the empty state
	return *state, nil
}

// Close implements port.StateRepository.
func (l *LocalStore) Close() {
}

// Store implements port.StateRepository.
func (l *LocalStore) Store(state model.State) error {
	s, _ := json.Marshal(state)
	return os.WriteFile(l.path+FILE, s, 0644)
}

func NewLocalStore(path string) port.StateRepository {
	return &LocalStore{path: path}
}
