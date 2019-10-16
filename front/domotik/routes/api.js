var express = require('express');
var router = express.Router();

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

router.get('/last/:time(24h|30d)/:sensor(meanPerHour|sumPerDay|tankHotWaterPerDay|livingRoomPerHour|waterPerDay)', function(req, res, next) {
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

              if (req.params.sensor == "tankHotWaterPerDay") {
                var _found = false;
                var date = new Date(new Date(entry.timestamp * 1000).toDateString()).getTime() / 1000;
                result[index].values.forEach(function(value) {
                  if(value[0] == date && !_found) {
                    value[1] += entry.value;
                    _found = true;
                  }
                });
                if (!_found) {
                  result[index].values.push([date, entry.value]);
                }
              } else {
                result[index].values.push([entry.timestamp, entry.value]);
              }
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
