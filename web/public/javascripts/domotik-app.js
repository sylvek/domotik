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
      $scope.twenty_four_hours_temp = [];
      $scope.twenty_four_hours_watt = [];
      $scope.twenty_four_hours_mean_watt = [];
      $scope.twenty_four_hours_presence_livingroom = [];
      $scope.thirty_days_sum_watt = [];
      $scope.thirty_days_hotwatertank_minutes = [];

      $scope.sum_watt_yesterday = "nc";
      $scope.sum_watt_last_year = "nc";
      $scope.sum_watt_now = "nc";
      $scope.mean_watt_last_30_days = "nc";
      $scope.hot_water_comsuption_yesterday = "nc";
      $scope.hot_water_mean_last_30_days = "nc";

      $scope.temp_bathroom = "nc";
      $scope.temp_outside = "nc";
      $scope.temp_living_room = "nc";
      $scope.power_hour = "nc";
      $scope.power_now = "nc";

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
          // bathroom
          case "esp8266":
            $scope.temp_bathroom = payload;
            break;
          // outside
          case "thn132n":
            $scope.temp_outside = payload;
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

        switch (categories[2]) {
          case "temp":
            var d = $scope.twenty_four_hours_temp;
            d.forEach(function(e2) {
              if (e2.key == categories[1]) {
                e2.values.push([new Date().getTime()/1000,payload]);
              }
            });
            $scope.twenty_four_hours_temp = d;
            break;
          case "watt":
            var d = $scope.twenty_four_hours_watt;
            d.forEach(function(e2) {
              if (e2.key == categories[1]) {
                e2.values.push([new Date().getTime()/1000,payload]);
              }
            });
            $scope.twenty_four_hours_watt = d;
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

      domotikSrv.last("24h", "temp").then(function(response) {
        $scope.twenty_four_hours_temp = response.data;
      });

      domotikSrv.last("24h", "watt").then(function(response) {
        $scope.twenty_four_hours_watt = response.data;
      });

      domotikSrv.last("24h", "meanPerHour").then(function(response) {
        $scope.twenty_four_hours_mean_watt = response.data;
        if (response.data.length > 0) {
          var sensor = response.data[0];
          var sum = 0;
          var count_this_date = false;
          sensor.values.forEach(function(e) {
            if (count_this_date) {
              sum += e[1];
            }
            var date = new Date(e[0] * 1000);
            if (date.getHours() == 0 && !count_this_date) {
              count_this_date = true;
            }
          });
          $scope.sum_watt_now = (sum / 1000).toPrecision(4);
        }
      });

      domotikSrv.last("24h", "livingRoomPerHour").then(function(response) {
        $scope.twenty_four_hours_presence_livingroom = response.data;
      });

      domotikSrv.last("30d", "sumPerDay").then(function(response) {
        $scope.thirty_days_sum_watt = response.data;
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
        $scope.thirty_days_hotwatertank_minutes = response.data;
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

      $scope.xAxisTickFormat_Date_Format = function() {
        return function(d) {
          return d3.time.format('%d/%m')(new Date(d * 1000));
        }
      }
      $scope.xAxisTickFormat_Time_Format = function() {
        return function(d) {
          return d3.time.format('%X')(new Date(d * 1000));
        }
      }
      $scope.xAxisTickFormat_Hour_Format = function() {
        return function(d) {
          return d3.time.format('%Hh')(new Date(d * 1000));
        }
      }
      $scope.yAxisTickFormat_Temp_Format = function() {
        return function(d) {
          return d + "°c";
        }
      }
      $scope.yAxisTickFormat_Watt_Format = function() {
        return function(d) {
          return d + "W";
        }
      }
      $scope.yAxisTickFormat_Min_Format = function() {
        return function(d) {
          return d + "min";
        }
      }
      $scope.yAxisTickFormat_Percent_Format = function() {
        return function(d) {
          return (d * 100).toPrecision(4) + "%";
        }
      }
  });
}(angular));
