package infrastructure

import (
	"encoding/json"
	"os"

	"github.com/sylvek/domotik/datastore/domain"
	"github.com/sylvek/domotik/datastore/port"
)

type LocalClient struct {
	path string
}

const FILE = "/state.json"

// Retrieve implements port.StateRepository.
func (l *LocalClient) Retrieve() (domain.State, error) {
	state := domain.NewState()
	data, err := os.ReadFile(l.path + FILE)
	if err == nil {
		json.Unmarshal(data, &state)
	}

	// if the file does not exist we return the empty state
	return *state, nil
}

// Close implements port.StateRepository.
func (l *LocalClient) Close() {
}

// Store implements port.StateRepository.
func (l *LocalClient) Store(state domain.State) error {
	s, _ := json.Marshal(state)
	return os.WriteFile(l.path+FILE, s, 0644)
}

func NewLocalClient(path string) port.StateRepository {
	return &LocalClient{path: path}
}
