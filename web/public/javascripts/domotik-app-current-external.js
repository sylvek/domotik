(function (angular) {
  var app = angular.module("domotikApp", ['nvd3ChartDirectives']);

  app.service('domotikSrv', function ($http) {

    return({
      last: last
    });

    function last(time, type) {
      return $http.get("api/last/" + time + "/" + type);
    }

  });

  app.controller('domotikIndexCtrl', function($scope, $interval, domotikSrv){

      $scope.sum_watt_yesterday = "--";
      $scope.sum_watt_last_year = "--";
      $scope.mean_watt_last_30_days = "--";
      $scope.hot_water_comsuption_yesterday = "--";
      $scope.hot_water_mean_last_30_days = "--";

      $scope.temp_outside1 = "--";
      $scope.temp_outside2 = "--";
      $scope.temp_bathroom = "--";
      $scope.temp_living_room = "--";
      $scope.power_hour = "--";
      $scope.power_now = "--";

      var client = new Paho.MQTT.Client("wss://sylvek.hd.free.fr/mosquitto", "gh-" + new Date().getTime());
      client.onConnectionLost = function(responseObject) {
          if (responseObject.errorCode !== 0) {
              console.log("onConnectionLost:" + responseObject.errorMessage);
              console.log("Reconnecting... [" + new Date() + "]");
              client.connect({
                  onSuccess: function() {
                      client.subscribe("sensors/#");
                  }
              });
          }
      };
      client.onMessageArrived = function(message) {
        var topic = message.destinationName;
        var payload = message.payloadString;
        console.log("onMessageArrived => " + topic + " " + payload);
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
          case "camera":
            $scope.last_view = "data:image/jpg;base64," + payload;
            break;
          default:
            break;
        }

        $scope.$apply();
      };
      client.connect({
        onSuccess: function() {
          console.log("onSuccess => subscribe to sensors/#");
          client.subscribe("sensors/#");
        }
      });

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
  });
}(angular));
