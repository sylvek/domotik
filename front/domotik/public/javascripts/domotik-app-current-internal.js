(function (angular) {
  var app = angular.module("domotikApp", ['nvd3ChartDirectives', 'ngtweet']);

  app.service('domotikSrv', function ($http) {

    return({
      last: last,
      getWeather: getWeather,
      quote: quote
    });

    function getWeather(query) {
      var weather = { temp: {}, clouds: null, humidity: null, pressure: null, wind: null };
      $http.jsonp('http://api.openweathermap.org/data/2.5/weather?q=' + query + '&units=metric&appid=5de8e795c11a94258c9e69cc8f2d87bf&callback=JSON_CALLBACK').success(function(data) {
          if (data) {
              if (data.main) {
                  weather.city = query;
                  weather.temp.current = data.main.temp;
                  weather.temp.min = data.main.temp_min;
                  weather.temp.max = data.main.temp_max;
                  weather.humidity = data.main.humidity;
                  weather.pressure = data.main.pressure;
                  weather.wind = data.wind.speed;
              }
              weather.clouds = data.clouds ? data.clouds.all : undefined;

              var baseUrl = 'https://ssl.gstatic.com/onebox/weather/128/';
                if (weather.clouds < 20) {
                    weather.img = baseUrl + 'sunny.png';
                } else if (weather.clouds < 90) {
                   weather.img = baseUrl + 'partly_cloudy.png';
                } else {
                    weather.img = baseUrl + 'cloudy.png';
                }
          }
      });

      return weather;
    }

    function last(time, type) {
      return $http.get("api/last/" + time + "/" + type);
    }

    function quote() {
      return $http.get("http://quotes.rest/qod.json");
    }

  });

  app.controller('domotikIndexCtrl', function($scope, $interval, domotikSrv){

      $scope.sum_watt_today = "--";
      $scope.sum_watt_yesterday = "--";
      $scope.sum_watt_last_year = "--";
      $scope.mean_watt_last_30_days = "--";
      $scope.hot_water_comsuption_yesterday = "--";
      $scope.hot_water_mean_last_30_days = "--";

      $scope.temp_outside = "--";
      $scope.temp_home = "--";
      $scope.power_hour = "--";
      $scope.power_now = "--";

      $scope.quote = "--";

      $scope.last_position_timestamp = 0;

      var client = new Paho.MQTT.Client(location.hostname, 9883, "gh-" + new Date().getTime());
      client.onConnectionLost = function(responseObject) {
          if (responseObject.errorCode !== 0) {
              console.log("onConnectionLost:" + responseObject.errorMessage);
              console.log("Reconnecting... [" + new Date() + "]");
              client.connect({
                  onSuccess: function() {
                      client.subscribe("current/+/temp");
                      client.subscribe("current/+/watt");
                      client.subscribe("sensors/+/temp");
                      client.subscribe("sensors/+/watt");
                  }
              });
          }
      };
      client.onMessageArrived = function(message) {
        var topic = message.destinationName;
        var payload = message.payloadString;
        var categories = topic.split("/");
        switch(categories[1]) {
          // room
          case "esp12e":
            $scope.temp_outside = payload;
            break;
          // home
          case "esp8266":
            $scope.temp_home = payload;
            break;
          // power
          case "linky":
            switch (categories[2]) {
              case "watt":
                $scope.power_now = (payload / 1000).toPrecision(4);
                break;
              default:
                break;
            }
            break;
          case "linkymean":
            $scope.power_hour = (payload / 1000).toPrecision(4);
            break;
          case "sumPerDay":
            $scope.sum_watt_today = (payload / 1000).toPrecision(4);
            break;
          default:
            break;
        }

        $scope.$apply();
      };
      client.connect({
        onSuccess: function() {
          console.log("onSuccess => subscribe to sensors, triggers (temp) & positions");
          client.subscribe("current/+/temp");
          client.subscribe("current/+/watt");
          client.subscribe("sensors/+/temp");
          client.subscribe("sensors/+/watt");
        }
      });

      var ephemeride = function() {
        $scope.ephemeride = ephemeris.getTodayEphemeris();
      }

      var quote = function() {
        domotikSrv.quote().then(function(response) {
          var quote = response.data.contents.quotes[0];
          $scope.quote = quote.quote + " - " + quote.author;
        });
      }

      var wallpaper = function() {
        random = Math.floor((Math.random() * 80) + 1);
        $scope.background = "img/" + random + ".jpg";
      }

      var twitter = function() {
        // twitter
        // https://github.com/arusahni/ngtweet/issues/30
        // took too much CPU.
        $scope.slug = "";
        var now = new Date();
        if (now.getDay() > 0 && now.getDay() < 6 && /* week days */
            now.getHours() > 5 && now.getHours() < 9 /* during departure */)
        {
          $scope.slug = "transilien";
        }
        if (now.getDay() > 0 && now.getDay() < 6 && /* week days */
            now.getHours() > 18 && now.getHours() < 23 /* during sunrise */)
        {
          $scope.slug = "news";
        }
      }

      var weather = function() {
        domotikSrv.last("30d", "sumPerDay").then(function(response) {
          if (response.data.length > 0) {
            var sensor = response.data[0];
            var number_of_values = sensor.values.length;
            if (number_of_values > 0) {
              $scope.sum_watt_yesterday = (sensor.values[number_of_values - 1][1] / 1000).toPrecision(4);

              var sum = 0;
              sensor.values.forEach(function(element) {
                sum += element[1];
              });
              $scope.mean_watt_last_30_days = (sum / number_of_values / 1000).toPrecision(4);
            }
          }
        });

        domotikSrv.last("30d", "tankHotWaterPerDay").then(function(response) {
          if (response.data.length > 0) {
            var sensor = response.data[0];
            var number_of_values = sensor.values.length;
            if (number_of_values > 0) {
              $scope.hot_water_comsuption_yesterday = (sensor.values[number_of_values - 1][1]).toPrecision(4);

              var sum = 0;
              sensor.values.forEach(function(element) {
                sum += element[1];
              });
              $scope.hot_water_mean_last_30_days = (sum / number_of_values).toPrecision(4);
            }
          }
        });

        domotikSrv.last("year", "sumPerDay").then(function(response) {
          if (response.data.length > 0) {
            $scope.sum_watt_last_year = (response.data[0].value / 1000).toPrecision(4);
          }
        });
      };

      weather();
      ephemeride();
      quote();
      wallpaper();
      twitter();

      // clock
      var timeController = this;
      timeController.clock = { time: "", interval: 1000 };
      $interval( function () {$scope.current_time = Date.now();}, timeController.clock.interval);

  });
}(angular));
