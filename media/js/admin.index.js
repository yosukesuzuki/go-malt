angular.module("myApp", [])
.controller("AdminController", function($scope, $http) {
    $http.get('/admin/rest/models').success(function(data) {
        $scope.models = data.models;
    });
});
