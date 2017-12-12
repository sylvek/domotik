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
      $scope.thirty_days_water_liters = [];

      domotikSrv.last("24h", "temp").then(function(response) {
        $scope.twenty_four_hours_temp = response.data;
      });

      domotikSrv.last("24h", "watt").then(function(response) {
        $scope.twenty_four_hours_watt = response.data;
      });

      domotikSrv.last("24h", "meanPerHour").then(function(response) {
        $scope.twenty_four_hours_mean_watt = response.data;
      });

      domotikSrv.last("24h", "livingRoomPerHour").then(function(response) {
        $scope.twenty_four_hours_presence_livingroom = response.data;
      });

      domotikSrv.last("30d", "sumPerDay").then(function(response) {
        $scope.thirty_days_sum_watt = response.data;
      });

      domotikSrv.last("30d", "tankHotWaterPerDay").then(function(response) {
        $scope.thirty_days_hotwatertank_minutes = response.data;
      });

      domotikSrv.last("30d", "waterPerDay").then(function(response) {
        $scope.thirty_days_water_liters = response.data;
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
      $scope.yAxisTickFormat_Liter_Format = function() {
        return function(d) {
          return d + "L";
        }
      }
  });
}(angular));
