package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sylvek/domotik/datastore/core"
	"github.com/sylvek/domotik/datastore/core/port"
	"github.com/sylvek/domotik/datastore/infrastructure"
)

func main() {
	serverPort := getenv("SERVER_PORT", "3333")
	basePath := strings.Trim(os.Getenv("BASE_PATH"), "/")
	mqttPort := getenv("MQTT_PORT", "1883")
	dbPath := os.Getenv("DB_PATH")
	dayEndedWebhookUrl := os.Getenv("DAY_ENDED_WEBHOOK_URL")
	dayEndedMeasures := strings.Fields(strings.ReplaceAll(getenv("DAY_ENDED_WEBHOOK_MEASURES", "daily_power_consumption"), ",", " "))
	settings := core.NewSettings(parseLabels(os.Getenv("SENSOR_LABELS")), dayEndedMeasures)

	log.Println("domotik starting...")

	broker, err := infrastructure.NewMQTTBroker(":" + mqttPort)
	if err != nil {
		log.Fatal(err)
	}
	localRepository := infrastructure.NewLocalStore(dbPath)
	sqliteRepository := infrastructure.NewSqliteDatabase(dbPath)

	application, err := core.NewApplication(localRepository, sqliteRepository, sqliteRepository, settings)
	if err != nil {
		log.Fatal(err)
	}

	sensors := broker.Messages("sensors/+/+")
	tasmota := broker.Messages("tele/+/SENSOR")
	if err := broker.Start(); err != nil {
		log.Fatal(err)
	}
	log.Printf("MQTT broker started (tcp: %s)", mqttPort)

	mqttPresenter := infrastructure.NewMqttPresenter(broker)
	go startSensorsListener(sensors, application, mqttPresenter)
	go startTasmotaListener(tasmota, application, mqttPresenter)
	go startHttpListener(serverPort, basePath, application)
	if dayEndedWebhookUrl != "" {
		go startDayEndedNotifier(application, infrastructure.NewWebhookNotifier(dayEndedWebhookUrl))
	}

	log.Println("domotik started")

	defer func() {
		broker.Stop()
		sqliteRepository.Close()
		localRepository.Close()
		log.Println("domotik stopped")
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// parseLabels reads a list like "living=séjour,room=chambre"
func parseLabels(value string) map[string]string {
	labels := map[string]string{}
	for _, pair := range strings.Split(value, ",") {
		if name, label, ok := strings.Cut(pair, "="); ok {
			labels[strings.TrimSpace(name)] = strings.TrimSpace(label)
		}
	}
	return labels
}

func startHttpListener(
	serverPort string,
	basePath string,
	application *core.Application) {

	dashboard := func(get func(port.DashboardPresenter) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := get(infrastructure.NewHttpDashboardPresenter(w)); err != nil {
				log.Printf(" - error - %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/summarize", func(w http.ResponseWriter, r *http.Request) {
		httpPresenter := infrastructure.NewHttpPresenter(w)
		application.GetSummarize(httpPresenter)
	})
	mux.HandleFunc("/api/today", dashboard(application.GetToday))
	mux.HandleFunc("/api/quarter", dashboard(application.GetQuarter))
	mux.HandleFunc("/api/history", dashboard(application.GetHistory))
	mux.Handle("/", infrastructure.NewDashboardHandler())

	// the application can be served behind a reverse proxy, with or without a path prefix
	var handler http.Handler = mux
	if basePath != "" {
		handler = http.NewServeMux()
		handler.(*http.ServeMux).Handle("/"+basePath+"/", http.StripPrefix("/"+basePath, mux))
		handler.(*http.ServeMux).Handle("/", mux)
	}

	log.Println("HTTP server started")

	err := http.ListenAndServe(":"+serverPort, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func startDayEndedNotifier(
	application *core.Application,
	notifier port.DayEndedNotifier) {

	log.Println("day ended webhook enabled")

	for range time.Tick(time.Minute) {
		if err := application.NotifyDayEnded(notifier); err != nil {
			log.Printf(" - error - %s", err)
		}
	}
}

func startTasmotaListener(
	messages chan infrastructure.Message,
	application *core.Application,
	publisher port.LogPublisher) {

	log.Println("tasmota converter started")

	for msg := range messages {
		telemetry, err := infrastructure.DecodeTasmota(msg)
		if err != nil {
			log.Printf(" - error - %s", err)
			continue
		}
		if err := application.ConvertTasmota(telemetry, publisher); err != nil {
			log.Printf(" - error - %s", err)
		}
	}
}

func startSensorsListener(
	messages chan infrastructure.Message,
	application *core.Application,
	presenter port.SummarizePresenter) {

	ticker := time.NewTicker(time.Minute)

	log.Println("sensors listener started")

	for {
		select {
		case <-ticker.C:
			if err := application.Process(presenter); err != nil {
				log.Fatal(err)
			}
		case msg := <-messages:
			elements := strings.Split(msg.Topic, "/")
			value, _ := strconv.ParseFloat(string(msg.Payload), 64)

			for {
				err := application.AddLog(elements, value)
				if err == nil {
					break
				}
				log.Printf(" - error - %s", err)
				time.Sleep(time.Second)
			}
		}
	}
}
