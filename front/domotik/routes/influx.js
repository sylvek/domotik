const Influx = require('influx')
const influx = new Influx.InfluxDB(process.env.INFLUX_DB || 'http://192.168.0.4:8086/domotik')

var express = require('express')
var router = express.Router()

// :type in { "daily_power_consumption", "daily_temp_outside", "daily_temp_inside" }
router.get('/last/24h/:type', function(req, res) {
    var query = `SELECT value as value FROM "infinite".` + req.params["type"] + ` WHERE time >= now() - 3d ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})
router.get('/last/30d/:type', function(req, res) {
    var query = `SELECT median(value) as value FROM "infinite".` + req.params["type"] + ` WHERE time >= now() - 30d ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})
router.get('/last/year/:type', function(req, res) {
    var query = `SELECT median(value) as value FROM "infinite".` + req.params["type"] + ` WHERE time > now() - 370d AND time < now() - 360d ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})

module.exports = router