package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sylvek/domotik/datastore/core/dto"
)

func TestWebhookNotifier(t *testing.T) {
	calls := make(chan *http.Request, 10)
	failures := 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if failures > 0 {
			failures--
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		calls <- r
	}))
	defer server.Close()

	notifier := NewWebhookNotifier(server.URL + "/exec?key=secret").(*WebhookNotifier)
	notifier.delay = time.Millisecond
	if err := notifier.Notify(dto.DayEnded{Timestamp: 1795132800, WattConsumed: 12345}); err != nil {
		t.Fatal(err)
	}

	// the first call fails, the second one succeeds
	select {
	case r := <-calls:
		query := r.URL.Query()
		if query.Get("title") != "💡 Consommation: 12.345 kWh" || query.Get("timestamp") != "1795132800" || query.Get("key") != "secret" {
			t.Errorf("unexpected query %v", query)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the webhook was not called")
	}
	select {
	case <-calls:
		t.Error("the webhook was called twice")
	case <-time.After(50 * time.Millisecond):
	}

	// without url, nothing is called
	for _, disabled := range []string{"", "not an url", "://"} {
		if err := NewWebhookNotifier(disabled).Notify(dto.DayEnded{}); err != nil {
			t.Error(err)
		}
	}
}
