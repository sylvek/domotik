package infrastructure

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/port"
)

type HttpPresenter struct {
	writer http.ResponseWriter
}

// Present implements port.MqttPresenter.
func (m *HttpPresenter) Present(output dto.Output) error {
	o, _ := json.Marshal(output)
	io.WriteString(m.writer, string(o))
	return nil
}

func NewHttpPresenter(w http.ResponseWriter) port.SummarizePresenter {
	return &HttpPresenter{writer: w}
}
