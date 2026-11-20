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

const (
	WEBHOOK_ATTEMPTS = 3
	WEBHOOK_TIMEOUT  = 30 * time.Second
)

// WebhookNotifier calls an URL (GET) with `title` and `timestamp` query parameters.
// Nothing is called when the URL is empty.
type WebhookNotifier struct {
	url    *url.URL
	delay  time.Duration
	client *http.Client
}

func title(event dto.DayEnded) string {
	return fmt.Sprintf("💡 Consommation: %s kWh", strconv.FormatFloat(float64(event.WattConsumed)/1000, 'f', -1, 64))
}

// Notify implements port.DayEndedNotifier.
// The call is done in background: a slow or unreachable endpoint never blocks the caller.
func (w *WebhookNotifier) Notify(event dto.DayEnded) error {
	if w.url == nil {
		return nil
	}
	target := *w.url
	query := target.Query()
	query.Set("title", title(event))
	query.Set("timestamp", strconv.FormatInt(event.Timestamp, 10))
	target.RawQuery = query.Encode()

	go func() {
		day := time.Unix(event.Timestamp, 0).UTC().Format("2006-01-02")
		for attempt := 1; attempt <= WEBHOOK_ATTEMPTS; attempt++ {
			body, err := w.call(target.String())
			if err == nil {
				log.Printf("webhook: %s notified - %s", day, body)
				return
			}
			log.Printf("webhook: %s not notified (attempt %d/%d) - %s", day, attempt, WEBHOOK_ATTEMPTS, err)
			if attempt < WEBHOOK_ATTEMPTS {
				time.Sleep(w.delay)
			}
		}
	}()

	return nil
}

func (w *WebhookNotifier) call(target string) (string, error) {
	response, err := w.client.Get(target)
	if err != nil {
		// the url is a secret, it must not be logged
		if urlError, ok := err.(*url.Error); ok {
			err = urlError.Err
		}
		return "", err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 200))

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status %d", response.StatusCode)
	}
	return string(body), nil
}

func NewWebhookNotifier(rawUrl string) port.DayEndedNotifier {
	notifier := &WebhookNotifier{delay: time.Minute, client: &http.Client{Timeout: WEBHOOK_TIMEOUT}}
	if rawUrl == "" {
		return notifier
	}

	target, err := url.Parse(rawUrl)
	if err != nil || target.Host == "" {
		// the url is a secret, it must not be logged
		log.Println("webhook: invalid url, the webhook is disabled")
		return notifier
	}
	notifier.url = target
	log.Println("webhook: enabled")
	return notifier
}
