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

  app.controller('domotikIndexCtrl', function($scope, domotikSrv){
      domotikSrv.last("24h", "temp").then(function(response) {
        $scope.twenty_four_hours_temp = response.data;
      });

      domotikSrv.last("24h", "watt").then(function(response) {
        $scope.twenty_four_hours_watt = response.data;
      });

      $scope.xAxisTickFormat_Date_Format = function() {
        return function(d) {
          return d3.time.format('%x')(new Date(d * 1000));
        }
      }
      $scope.xAxisTickFormat_Time_Format = function() {
        return function(d) {
          return d3.time.format('%X')(new Date(d * 1000));
        }
      }
      $scope.yAxisTickFormat_Temp_Format = function() {
        return function(d) {
          return d + "°c";
        }
      }
      $scope.yAxisTickFormat_Watt_Format = function() {
        return function(d) {
          return d + " watt";
        }
      }
  });
}(angular));
