(function (angular) {
  var app = angular.module('domotikApp', ['ngtweet']);
  var magical = {
    0: 90000,
    1: 60000,
    2: 50000,
    3: 47000,
    4: 45000,
    5: 44000,
    6: 43000,
    7: 42500,
    8: 37000,
    9: 34000,
    10: 34000,
    11: 28000,
    12: 26500,
    13: 23500,
    14: 23000,
    15: 19500,
    16: 17500,
    17: 17500,
    18: 17000,
    19: 16500,
    20: 15000,
    21: 14000,
    22: 14000,
    23: 14000,
    24: 13000,
    25: 12000,
    26: 11500,
    27: 11500,
    28: 11000,
    29: 10000,
    30: 9000,
  };

  app.service('domotikSrv', function ($http) {
    return {
      last: last,
      quote: quote,
    };

    function last(time, type) {
      return $http.get('api/last/' + time + '/' + type);
    }

    function quote() {
      return $http.get('http://quotes.rest/qod.json');
    }
  });

  app.controller(
    'domotikIndexCtrl',
    function ($scope, $interval, $window, domotikSrv) {
      $scope.sum_watt_today = '--';
      $scope.sum_euro_today = '--';
      $scope.sum_watt_yesterday = '--';
      $scope.sum_watt_last_year = '--';
      $scope.mean_watt_last_30_days = '--';

      $scope.outside_temp_yesterday = '--';
      $scope.outside_temp_last_year = '--';
      $scope.outside_temp_last_30_days = '--';

      $scope.consumption_scoring_yesterday = '--';
      $scope.consumption_scoring_last_year = '--';
      $scope.consumption_scoring_last_30_days = '--';

      $scope.temp_outside = '--';
      $scope.temp_living = '--';
      $scope.temp_room = '--';
      $scope.power_hour = '--';
      $scope.power_now = '--';

      $scope.quote = '--';

      $scope.last_position_timestamp = 0;

      var client = new Paho.MQTT.Client(
        location.hostname,
        9883,
        'gh-' + new Date().getTime()
      );
      client.onConnectionLost = function (responseObject) {
        if (responseObject.errorCode !== 0) {
          console.log('onConnectionLost:' + responseObject.errorMessage);
          console.log('Reconnecting... [' + new Date() + ']');
          client.connect({
            onSuccess: function () {
              client.subscribe('sensors/+/temp');
              client.subscribe('sensors/+/watt');
              client.subscribe('sensors/+/euro');
            },
          });
        }
      };
      client.onMessageArrived = function (message) {
        var topic = message.destinationName;
        var payload = message.payloadString;
        var categories = topic.split('/');
        switch (categories[1]) {
          // room
          case 'outside':
            $scope.temp_outside = payload;
            break;
          // home
          case 'living':
            $scope.temp_living = payload;
            break;
          case 'room':
            $scope.temp_room = payload;
            break;
          // power
          case 'linky':
            switch (categories[2]) {
              case 'watt':
                $scope.power_now = (payload / 1000).toFixed(2);
                break;
              default:
                break;
            }
            break;
          case 'meanPerHour':
            switch (categories[2]) {
              case 'watt':
                $scope.power_hour = (payload / 1000).toFixed(2);
                break;
              default:
                break;
            }
          case 'sumPerDay':
            switch (categories[2]) {
              case 'watt':
                $scope.sum_watt_today = (payload / 1000).toFixed(2);
                break;
              case 'euro':
                $scope.sum_euro_today = payload;
                break;
              default:
                break;
            }
            break;
          default:
            break;
        }

        $scope.$apply();
      };
      client.connect({
        onSuccess: function () {
          console.log(
            'onSuccess => subscribe to sensors, triggers (temp) & positions'
          );
          client.subscribe('sensors/+/temp');
          client.subscribe('sensors/+/watt');
          client.subscribe('sensors/+/euro');
        },
      });

      var ephemeride = function () {
        $scope.ephemeride = ephemeris.getTodayEphemeris();
      };

      var quote = function () {
        domotikSrv.quote().then(function (response) {
          var quote = response.data.contents.quotes[0];
          $scope.quote = quote.quote + ' - ' + quote.author;
        });
      };

      var wallpaper = function () {
        var keys = Object.keys($window.images);
        random = Math.floor(Math.random() * keys.length);
        var wallpaper = $window.images[keys[random]];
        $scope.background = 'img/' + wallpaper.file;
        $scope.title = wallpaper.title;
      };

      var twitter = function () {
        // twitter
        // https://github.com/arusahni/ngtweet/issues/30
        // took too much CPU.
        $scope.slug = '';
        var now = new Date();
        if (
          now.getDay() > 0 &&
          now.getDay() < 6 /* week days */ &&
          now.getHours() > 5 &&
          now.getHours() < 9 /* during departure */
        ) {
          $scope.slug = 'transilien';
        }
        if (
          now.getDay() > 0 &&
          now.getDay() < 6 /* week days */ &&
          now.getHours() > 18 &&
          now.getHours() < 23 /* during sunrise */
        ) {
          $scope.slug = 'news';
        }
      };

      var scoring = function (power, temperature) {
        if (!isNaN(power) && !isNaN(temperature)) {
          let score = (
            (power * 1000) /
            magical[Math.floor(temperature)]
          ).toFixed(2);
          if (score > 1.15) {
            return 'red';
          }
          if (score < 1.0) {
            return 'green';
          }
        }
        return 'yellow';
      };

      var sums = function () {
        domotikSrv
          .last('24h', 'daily_power_consumption')
          .then(function (response) {
            if (response.data.length > 0) {
              $scope.sum_watt_yesterday = (
                response.data[0].value / 1000
              ).toFixed(2);
              $scope.consumption_scoring_yesterday = scoring(
                $scope.sum_watt_yesterday,
                $scope.outside_temp_yesterday
              );
            }
          });

        domotikSrv
          .last('30d', 'daily_power_consumption')
          .then(function (response) {
            if (response.data.length > 0) {
              $scope.mean_watt_last_30_days = (
                response.data[0].value / 1000
              ).toFixed(2);
              $scope.consumption_scoring_last_30_days = scoring(
                $scope.mean_watt_last_30_days,
                $scope.outside_temp_last_30_days
              );
            }
          });

        domotikSrv
          .last('year', 'daily_power_consumption')
          .then(function (response) {
            if (response.data.length > 0) {
              $scope.sum_watt_last_year = (
                response.data[0].value / 1000
              ).toFixed(2);
              $scope.consumption_scoring_last_year = scoring(
                $scope.sum_watt_last_year,
                $scope.outside_temp_last_year
              );
            }
          });
      };

      var temperatures = function () {
        domotikSrv.last('24h', 'daily_temp_outside').then(function (response) {
          if (response.data.length > 0) {
            $scope.outside_temp_yesterday = response.data[0].value.toFixed(2);
            $scope.consumption_scoring_yesterday = scoring(
              $scope.sum_watt_yesterday,
              $scope.outside_temp_yesterday
            );
          }
        });

        domotikSrv.last('30d', 'daily_temp_outside').then(function (response) {
          if (response.data.length > 0) {
            $scope.outside_temp_last_30_days =
              response.data[0].value.toFixed(2);
            $scope.consumption_scoring_last_30_days = scoring(
              $scope.mean_watt_last_30_days,
              $scope.outside_temp_last_30_days
            );
          }
        });

        domotikSrv.last('year', 'daily_temp_outside').then(function (response) {
          if (response.data.length > 0) {
            $scope.outside_temp_last_year = response.data[0].value.toFixed(2);
            $scope.consumption_scoring_last_year = scoring(
              $scope.sum_watt_last_year,
              $scope.outside_temp_last_year
            );
          }
        });
      };

      sums();
      temperatures();
      ephemeride();
      quote();
      wallpaper();
      twitter();

      // clock
      var timeController = this;
      timeController.clock = { time: '', interval: 1000 };
      $interval(function () {
        $scope.current_time = Date.now();
      }, timeController.clock.interval);
    }
  );
})(angular);
