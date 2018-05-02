// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$scope.pets = [];

	$scope.init = function(){
		var petsRow = document.getElementById('petsRow');
		var petTemplate = document.getElementById('petTemplate');
		  
		$scope.getAllPets(function(pets) {
			pets.forEach((pet) => {
				var template = petTemplate.cloneNode(true);
				template.style.display = "block";

				template.querySelector('.panel-title').innerHTML = pet.name;
				template.querySelector('img').setAttribute('src', pet.picture)
				template.querySelector('.pet-breed').innerHTML = pet.breed;
				template.querySelector('.pet-location').innerHTML = pet.location;
				template.querySelector('.pet-age').innerHTML = pet.age;

				var button = template.querySelector('.btn-adopt')
				button.setAttribute('data-id', pet.Key)
				button.addEventListener('click', $scope.handleAdoption);

				if (pet.owner) {
					$scope.markAdopted(button);
				}

				petsRow.appendChild(template);			
			});
		});
	};

	$scope.handleAdoption = function(event){
		event.preventDefault();

		var petId = parseInt(event.target.getAttribute('data-id'));

		appFactory.adopt(petId, function(data){
			$scope.markAdopted(event.target);
		});
	};

	$scope.markAdopted = function(button){
		button.innerHTML = "Success";
		button.disabled = true;
	};
	
	$scope.getAllPets = function(callback){
		appFactory.queryAllPets(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
				return parseFloat(a.Key) - parseFloat(b.Key);
			});
			
			callback(array);
		});
	};
});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllPets = function(callback){

    	$http.get('/pets/').success(function(output){
			callback(output)
		});
	}

	factory.adopt = function(petId, callback){

    	$http.get('/adoptPet/'+petId).success(function(output){
			callback(output)
		});
	}

	return factory;
});
