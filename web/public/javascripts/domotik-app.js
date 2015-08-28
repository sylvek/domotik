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
      $scope.thirty_days_sum_watt = [];
      $scope.lastCapture = "api/last/capture";

      $scope.sum_watt_yesterday = "nc";
      $scope.sum_watt_last_year = "nc";
      $scope.sum_watt_now = "nc";
      $scope.mean_last_30_days = "nc";

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
            $scope.mean_last_30_days = (sum / number_of_values / 1000).toPrecision(4);
          }
        }
      });

      domotikSrv.last("year", "sumPerDay").then(function(response) {
        if (response.data.length > 0) {
          $scope.sum_watt_last_year = (response.data[0].value / 1000).toPrecision(4);
        }
      });

      function update() {
        $scope.lastCapture = "api/last/capture?" + new Date().getTime();
        domotikSrv.last("20s", "temp").then(function(response) {
          response.data.forEach(function(e1) {
            var d = $scope.twenty_four_hours_temp;
            d.forEach(function(e2) {
                if (e2.key == e1.key) {
                  e1.values.forEach(function(e3) {
                    e2.values.push(e3);
                  });
                }
            });
            $scope.twenty_four_hours_temp = d;
          });
        });
        domotikSrv.last("20s", "watt").then(function(response) {
          response.data.forEach(function(e1) {
            var d = $scope.twenty_four_hours_watt;
            d.forEach(function(e2) {
                if (e2.key == e1.key) {
                  var last = 0;
                  e1.values.forEach(function(e3) {
                    e2.values.push(e3);
                  });
                }
            });
            $scope.twenty_four_hours_watt = d;
          });
        })
      }
      setInterval(function() {$scope.$apply(update);}, 20000 /* 20S */);

      $scope.xAxisTickFormat_Date_Format = function() {
        return function(d) {
          return d3.time.format('%e/%m')(new Date(d * 1000));
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
  });
}(angular));
