const Influx = require('influx')
const influx = new Influx.InfluxDB(process.env.INFLUX_DB || 'http://192.168.0.4:8086/domotik')

var express = require('express')
var router = express.Router()

router.get('/last/24h/sumPerDay', function(req, res) {
    var query = `SELECT median("value") as value FROM "measures" WHERE ("name" = 'sumPerDay' AND time >= now() - 24h) GROUP BY time(24h) fill(none) ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})
router.get('/last/30d/sumPerDay', function(req, res) {
    var query = `SELECT median("value") as value FROM "infinite"."daily_power_consumption" WHERE time >= now() - 30d GROUP BY time(30d) fill(none) ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})
router.get('/last/year/sumPerDay', function(req, res) {
    var query = `SELECT median("value") as value FROM "infinite"."daily_power_consumption" WHERE time > now() - 370d AND time < now() - 360d GROUP BY time(7d) fill(none) ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})

module.exports = router