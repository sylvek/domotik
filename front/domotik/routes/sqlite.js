const sqlite3 = require('sqlite3');
const db = new sqlite3.Database(process.env.DB_PATH + '/history.db', sqlite3.OPEN_READONLY);

var express = require('express')
var router = express.Router()

// :type in { "daily_power_consumption", "daily_temp_outside", "daily_temp_inside" }
router.get('/last/24h/:type', function(req, res) {
    var query = `SELECT value as value FROM data WHERE name="` + req.params["type"] + `" AND ts>=strftime("%s", "now", "-3 day") ORDER BY ts DESC LIMIT 1`;
    db.all(query, (err, rows) => { 
        if (err) {
            res.status(500).send(err.stack)
        }
        res.json(rows)
    })
})
router.get('/last/30d/:type', function(req, res) {
    var query = `SELECT avg(value) as value FROM data WHERE name="` + req.params["type"] + `" AND ts>=strftime("%s", "now", "-30 day") ORDER BY ts DESC LIMIT 1`;
    db.all(query, (err, rows) => { 
        if (err) {
            res.status(500).send(err.stack)
        }
        res.json(rows)
    })
})
router.get('/last/year/:type', function(req, res) {
    var query = `SELECT avg(value) as value FROM data WHERE name="` + req.params["type"] + `" AND ts>strftime("%s", "now", "-370 day") AND ts<strftime("%s", "now", "-360 day") ORDER BY ts DESC LIMIT 1`;
    db.all(query, (err, rows) => { 
        if (err) {
            res.status(500).send(err.stack)
        }
        res.json(rows)
    })
})

module.exports = router