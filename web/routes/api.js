var express = require('express');
var router = express.Router();

router.post('/google-home-request', function(req, res, next) {
  var db = req.db;
  var result = {};

  res.set('Content-Type', 'application/json');

  var type = req.body.result.parameters.type;
  var table = ("watt" == type) ? "measures" : "sensors";
  var sensor = ("watt" == type) ? "sumPerDay" : req.body.result.parameters.sensors;
  var unit = ("watt" == type) ? " watt" : " degrés";

  var collection = db.get(table);
  collection.ensureIndex({ type:1, timestamp:1 }, { background:true }, function(err, indexName) {
    collection.find({ type: type, sensor: sensor },
                    { fields: {value:1, _id:0},
                      sort: {timestamp:-1},
                      limit: 1},
                      function (err, elements) {
          result.speech = "aucune valeur trouvée";
          result.displayText = "aucune valeur trouvée";

          elements.forEach(function(entry){
                result.speech = entry.value + unit;
                result.displayText = entry.value + unit;
          });

          res.send(result);
    });
  });
});

router.get('/last/:time(24h|20s)/:type(temp|watt)', function(req, res, next) {
  var db = req.db;
  var collection = db.get("sensors");
  var twenty_four_hours = Math.floor((new Date().getTime() - 60*60*24*1000) / 1000);
  var twenty_seconds = Math.floor((new Date().getTime() - 20*1000) / 1000);
  var timestamp = ("24h" == req.params.time) ? twenty_four_hours : twenty_seconds;

  collection.ensureIndex({ type:1, timestamp:1 }, { background:true }, function(err, indexName) {
    collection.find({ type:req.params.type, timestamp:{$gt: timestamp} },
                    { fields: {timestamp:1, sensor:1, value:1, _id:0},
                      sort: {timestamp:1} }, function (err, elements) {

          var result = [];
          var type = [];
          elements.forEach(function(entry){
              var index = type.indexOf(entry.sensor);
              if (index == -1) {
                type.push(entry.sensor);
                result.push({"key": entry.sensor, "values":[]});
                index = type.length - 1;
              }
              result[index].values.push([entry.timestamp, entry.value]);
          });

          res.send(result);
    });
  });
});

router.get('/last/:time(24h|30d)/:sensor(meanPerHour|sumPerDay|tankHotWaterPerDay|livingRoomPerHour)', function(req, res, next) {
  var db = req.db;
  var collection = db.get("measures");
  var twenty_four_hours = Math.floor((new Date().getTime() - 60*60*24*1000) / 1000);
  var thirty_days = Math.floor((new Date().getTime() - 30*60*60*24*1000) / 1000);
  var timestamp = ("24h" == req.params.time) ? twenty_four_hours : thirty_days;

  collection.ensureIndex({ sensors:1, type:1, timestamp:1 }, { background:true }, function(err, indexName) {
    collection.find({ sensor:req.params.sensor, timestamp:{$gt: timestamp} },
                    { fields: {timestamp:1, sensor:1, value:1, _id:0},
                      sort: {timestamp:1} }, function (err, elements) {

          var result = [];
          var type = [];
          elements.forEach(function(entry){
              var index = type.indexOf(entry.sensor);
              if (index == -1) {
                type.push(entry.sensor);
                result.push({"key": entry.sensor, "values":[]});
                index = type.length - 1;
              }
              result[index].values.push([entry.timestamp, entry.value]);
          });

          res.send(result);
    });
  });
});

router.get('/last/year/:sensor(meanPerHour|sumPerDay)', function(req, res, next) {
  var db = req.db;
  var collection = db.get("measures");

  var now = new Date();
  var last_year_date = new Date();
  last_year_date.setFullYear(now.getFullYear() - 1);
  var last_year = Math.floor(last_year_date.getTime() / 1000);

  collection.ensureIndex({ sensors:1, type:1, timestamp:1 }, { background:true }, function(err, indexName) {
    collection.find({ sensor:req.params.sensor, type:'watt', timestamp:{$gt: last_year} },
                    { fields: {timestamp:1, sensor:1, value:1, _id:0},
                      sort: {timestamp:1}, limit:1 }, function (err, elements) {
          res.send(elements);
    });
  });
});

module.exports = router;
