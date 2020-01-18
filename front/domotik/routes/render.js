var express = require('express');
var router = express.Router();

/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('current-internal', { title: 'domotik' });
});
router.get('/internal', function(req, res, next) {
  res.render('current-internal', { title: 'domotik' });
});

module.exports = router;
