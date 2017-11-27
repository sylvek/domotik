(function (angular) {
  var app = angular.module("domotikApp", ['nvd3ChartDirectives', 'ngtweet']);

  app.service('domotikSrv', function ($http) {

    return({
      last: last,
      getWeather: getWeather
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

  });

  app.controller('domotikIndexCtrl', function($scope, $interval, domotikSrv){

      $scope.sum_watt_today = "--";
      $scope.sum_watt_yesterday = "--";
      $scope.sum_watt_last_year = "--";
      $scope.mean_watt_last_30_days = "--";
      $scope.hot_water_comsuption_yesterday = "--";
      $scope.hot_water_mean_last_30_days = "--";

      $scope.temp_outside1 = "--";
      $scope.moisture_outside1 = "--";
      $scope.battery_outside1 = "--";
      $scope.temp_outside2 = "--";
      $scope.temp_bathroom = "--";
      $scope.temp_living_room = "--";
      $scope.power_hour = "--";
      $scope.power_now = "--";

      $scope.last_position_timestamp = 0;

      var client = new Paho.MQTT.Client("ws://192.168.0.2/mosquitto", "gh-" + new Date().getTime());
      client.onConnectionLost = function(responseObject) {
          if (responseObject.errorCode !== 0) {
              console.log("onConnectionLost:" + responseObject.errorMessage);
              console.log("Reconnecting... [" + new Date() + "]");
              client.connect({
                  onSuccess: function() {
                      client.subscribe("current/+/temp");
                      client.subscribe("current/+/watt");
                      client.subscribe("sharemyposition/+/position");
                      client.subscribe("sensors/+/temp");
                      client.subscribe("sensors/+/watt");
                      client.subscribe("sensors/+/battery");
                      client.subscribe("sensors/+/moisture");
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
            $scope.temp_outside2 = payload;
            break;
          // bathroom
          case "esp8266":
            $scope.temp_bathroom = payload;
            break;
          // outside
          case "thn132n":
            switch (categories[2]) {
              case "temp":
                $scope.temp_outside1 = payload;
                break;
              case "moisture":
                $scope.moisture_outside1 = payload;
                break;
              case "battery":
                $scope.battery_outside1 = payload;
                break;
              default:
                break;
            }
            break;
          // living room
          case "cc128":
            switch (categories[2]) {
              case "temp":
                $scope.temp_living_room = payload;
                break;
              case "watt":
                $scope.power_now = (payload / 1000).toPrecision(4);
                break;
              default:
                break;
            }
            break;
          case "cc128mean":
            $scope.power_hour = (payload / 1000).toPrecision(4);
            break;
          case "sumPerDay":
            $scope.sum_watt_today = (payload / 1000).toPrecision(4);
            break;
          default:
            switch (categories[2]) {
              case "position":
                var p = JSON.parse(payload);
                if ($scope.last_position_timestamp < p.timestamp) {
                  var position = p.lat + "," + p.lng;
                  var newDate = new Date();
                  newDate.setTime(p.timestamp);
                  $scope.last_date = newDate.toLocaleString();
                  $scope.last_position_timestamp = p.timestamp;
                  $scope.last_position = "http://staticmap.openstreetmap.de/staticmap.php?center=" + position + "&zoom=12&size=420x600&maptype=mapnik&markers=" + position + ",ol-marker-gold"
                }
                break;
              default:
                break;
            }
            break;
        }

        $scope.$apply();
      };
      client.connect({
        onSuccess: function() {
          console.log("onSuccess => subscribe to sensors, triggers (temp) & positions");
          client.subscribe("current/+/temp");
          client.subscribe("current/+/watt");
          client.subscribe("sharemyposition/+/position");
          client.subscribe("sensors/+/temp");
          client.subscribe("sensors/+/watt");
          client.subscribe("sensors/+/battery");
          client.subscribe("sensors/+/moisture");
        }
      });

      var ephemeride = function() {
        $scope.ephemeride = ephemeris.getTodayEphemeris();
      }

      var weather = function() {
        //$scope.weather1 = domotikSrv.getWeather("Juvisy,FR");
        //$scope.weather2 = domotikSrv.getWeather("Santa Monica,CA");
        //$scope.weather3 = domotikSrv.getWeather("New York,NY");

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
          now.getHours() > 17 && now.getHours() < 23 /* during sunrise */)
      {
        $scope.slug = "news";
      }

      // wallpaper
      switch (now.getMonth()) {
        default:
        case 0:
          $scope.background = "images/newyear.jpg";
          break;
        case 1:
        case 2:
          $scope.background = "images/winter.jpg";
          break;
        case 3:
        case 4:
          $scope.background = "images/motherday.jpg";
          break;
        case 5:
          $scope.background = "images/spring.jpg";
          break;
        case 6:
        case 7:
        case 8:
          $scope.background = "images/summer.jpg";
          break;
        case 9:
        case 10:
          $scope.background = "images/autumn.jpg";
          break;
        case 11:
          $scope.background = "images/christmas.jpg";
          break;
      }

      // clock
      var timeController = this;
      timeController.clock = { time: "", interval: 1000 };
      $interval( function () {$scope.current_time = Date.now();}, timeController.clock.interval);

  });
}(angular));
