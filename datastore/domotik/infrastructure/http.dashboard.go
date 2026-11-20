package infrastructure

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/port"
)

//go:embed web
var web embed.FS

// NewDashboardHandler serves the web pages displaying the dashboards.
func NewDashboardHandler() http.Handler {
	content, err := fs.Sub(web, "web")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(content))
}

type HttpDashboardPresenter struct {
	writer http.ResponseWriter
}

func (h *HttpDashboardPresenter) present(v any) error {
	h.writer.Header().Set("Content-Type", "application/json")
	h.writer.Header().Set("Cache-Control", "no-store")
	return json.NewEncoder(h.writer).Encode(v)
}

// PresentToday implements port.DashboardPresenter.
func (h *HttpDashboardPresenter) PresentToday(today dto.Today) error {
	return h.present(today)
}

// PresentQuarter implements port.DashboardPresenter.
func (h *HttpDashboardPresenter) PresentQuarter(quarter dto.Quarter) error {
	return h.present(quarter)
}

// PresentHistory implements port.DashboardPresenter.
func (h *HttpDashboardPresenter) PresentHistory(history dto.History) error {
	return h.present(history)
}

func NewHttpDashboardPresenter(w http.ResponseWriter) port.DashboardPresenter {
	return &HttpDashboardPresenter{writer: w}
}
