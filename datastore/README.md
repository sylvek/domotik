# domotik-datastore

One single binary that:

- embeds a MQTT broker (TCP), no need to run mosquitto,
- stores every `sensors/+/+` message into SQLite,
- converts Tasmota telemetry (`tele/{device}/SENSOR`) into `sensors/{device}/temp|humidity|dewpoint`,
- serves the dashboards (today / quarter / history),
- calls a webhook when a day ends (optional).

## build it

```
$> docker build -t domotik-datastore .
```

## run it

```
$> docker run -d --name datastore -p 3333:3333 -p 1883:1883 -e TZ=Europe/Paris -v /path/to/database:/database domotik-datastore
```

## configure it

| variable            | default      | usage                                                                  |
| ------------------- | ------------ | ---------------------------------------------------------------------- |
| SERVER_PORT         | 3333         | HTTP port (dashboards, `/api/*`, `/summarize`)                         |
| BASE_PATH           |              | path prefix when the reverse proxy does not strip it (ie. `/domotik`)  |
| MQTT_PORT           | 1883         | MQTT broker port                                                       |
| DB_PATH             | /database/   | where `sensors.db`, `history.db` and `state.json` are stored           |
| SENSOR_LABELS       |              | names displayed on dashboards, ie. `living=séjour,room=chambre`        |
| DAY_ENDED_WEBHOOK_URL      |       | URL called when a day ends, disabled when empty                 |
| TZ                  |              | timezone used to split days, ie. `Europe/Paris`                        |

## day ended webhook

When a new day starts, the measures of the day that has just ended are aggregated into `history.db`.
If `DAY_ENDED_WEBHOOK_URL` is set, the URL is then called (GET) with the electricity consumed during that day, as two query parameters:

| parameter | example                       | usage                                   |
| --------- | ----------------------------- | --------------------------------------- |
| title     | `💡 Consommation: 12.345 kWh` | the consumption of the day, in kWh      |
| timestamp | `1795132800`                  | the ended day, at midnight UTC, seconds |

A failed call is retried twice, one minute apart. The URL is a secret: it is given by the environment only and is never logged.

```
$> docker run -d --name datastore -e DAY_ENDED_WEBHOOK_URL=https://script.google.com/macros/s/xxx/exec …
```

## behind nginx

Every URL used by the dashboards is relative, so it can be served under any path.

```
location /domotik/ {
    proxy_pass http://127.0.0.1:3333/;
}
```

## architecture

```
core/                 business rules, no dependency
  model/              entities (log, measure, state, point, tasmota telemetry)
  port/               interfaces required by the core (repositories, presenters, publisher)
  dto/                data given to presenters
  application.go      use cases: store logs, summarize, convert tasmota
  dashboard.go        use cases: today, quarter, history
infrastructure/       adapters (sqlite, embedded mqtt broker, http, web pages, webhook)
main.go               wiring
```
