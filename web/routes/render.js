var express = require('express');
var router = express.Router();

/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('index', { title: 'domotik' });
});
router.get('/current', function(req, res, next) {
  res.render('current', { title: 'domotik' });
});
router.get('/history', function(req, res, next) {
  res.render('history', { title: 'domotik' });
});

module.exports = router;
