package infrastructure

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/port"
)

// WebhookNotifier calls an URL (GET) with `title` and `timestamp` query parameters.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

func title(event dto.DayEnded) string {
	return fmt.Sprintf("💡 Consommation: %s kWh", strconv.FormatFloat(float64(event.WattConsumed)/1000, 'f', -1, 64))
}

// Notify implements port.DayEndedNotifier.
func (w *WebhookNotifier) Notify(event dto.DayEnded) error {
	target, err := url.Parse(w.url)
	if err != nil {
		return fmt.Errorf("webhook: invalid url")
	}
	query := target.Query()
	query.Set("title", title(event))
	query.Set("timestamp", strconv.FormatInt(event.Timestamp, 10))
	target.RawQuery = query.Encode()

	response, err := w.client.Get(target.String())
	if err != nil {
		// the url is a secret, it must not be logged
		if urlError, ok := err.(*url.Error); ok {
			err = urlError.Err
		}
		return fmt.Errorf("webhook: %w", err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 200))

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", response.StatusCode)
	}

	log.Printf("webhook: %s notified - %s", time.Unix(event.Timestamp, 0).UTC().Format("2006-01-02"), body)
	return nil
}

func NewWebhookNotifier(url string) port.DayEndedNotifier {
	return &WebhookNotifier{url: url, client: &http.Client{Timeout: 30 * time.Second}}
}
