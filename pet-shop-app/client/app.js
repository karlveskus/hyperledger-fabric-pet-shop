// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){
	
	$scope.pets = appFactory.queryAllPets(function(data){
		var array = [];
		for (var i = 0; i < data.length; i++){
			parseInt(data[i].Key);
			data[i].Record.Key = parseInt(data[i].Key);
			array.push(data[i].Record);
		}
		array.sort(function(a, b) {
			return parseFloat(a.Key) - parseFloat(b.Key);
		});
		$scope.pets = array;
	});
	
});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllPets = function(callback){

    	$http.get('/pets/').success(function(output){
			callback(output)
		});
	}

	return factory;
});
