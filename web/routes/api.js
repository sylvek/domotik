var express = require('express');
var router = express.Router();

router.get('/last/:time(24h|20s)/:type(temp|watt)', function(req, res, next) {
  var db = req.db;
  var collection = db.get("sensors");
  var twenty_four_hours = Math.floor((new Date().getTime() - 60*60*24*1000) / 1000);
  var twenty_seconds = Math.floor((new Date().getTime() - 20*1000) / 1000);
  var timestamp = ("24h" == req.params.time) ? twenty_four_hours : twenty_seconds;

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

module.exports = router;
