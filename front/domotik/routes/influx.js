const Influx = require('influx')
const influx = new Influx.InfluxDB(process.env.INFLUX_DB || 'http://192.168.0.4:8086/domotik')

var express = require('express')
var router = express.Router()

router.get('/last/30d/:sensor(meanPerHour|sumPerDay)', function(req, res) {
    var query = `SELECT median("value") as value FROM "measures" WHERE ("name" = '${req.params.sensor}' AND time >= now() - 30d) GROUP BY time(30d) ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})
router.get('/last/24h/:sensor(meanPerHour|sumPerDay)', function(req, res) {
    var query = `SELECT last("value") as value FROM "measures" WHERE ("name" = '${req.params.sensor}' AND time >= now() - 2d) GROUP BY time(1d) LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})
router.get('/last/year/:sensor(meanPerHour|sumPerDay)', function(req, res) {
    var query = `SELECT median("value") as value FROM "measures" WHERE ("name" = '${req.params.sensor}') AND time > now() - 370d AND time < now() - 360d GROUP BY time(7d) ORDER BY time DESC LIMIT 1`;
    influx
        .query(query)
        .then(result => { res.json(result) })
        .catch(err => { res.status(500).send(err.stack) })
})

module.exports = router