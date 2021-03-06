//SPDX-License-Identifier: Apache-2.0

var pets = require('./controller.js');

module.exports = function(app){
  app.get('/pets', function(req, res){
    pets.getAllPets(req, res);
  });
  app.get('/adoptPet/:id', function(req, res){
    pets.adoptPet(req, res);
  })
}
