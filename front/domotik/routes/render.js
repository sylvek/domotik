var express = require('express');
var router = express.Router();

/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('index', { title: 'domotik' });
});
router.get('/internal', function(req, res, next) {
  res.render('current-internal', { title: 'domotik' });
});
router.get('/external', function(req, res, next) {
  res.render('current-external', { title: 'domotik' });
});
router.get('/history', function(req, res, next) {
  res.render('history', { title: 'domotik' });
});

module.exports = router;
