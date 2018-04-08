var mongoose = require('mongoose');
var express = require('express');
var express_graphql = require('express-graphql');
var { buildSchema } = require('graphql');

// mongoose connection
mongoose.connect(process.argv[2], { useMongoClient: true });
mongoose.Promise = global.Promise;

var db = mongoose.connection;
db.on('error', console.error.bind(console, 'connection error:'));
db.once('open', function() {

  // mongoose model
  var MeasureModel = mongoose.model('Measure', {
    type: String,
    sensor: String,
    value: Number,
    timestamp: Number
  }, 'measures');

  // GraphQL schema
  var schema = buildSchema(`
    type Query {
      measures(sensor: String!): [Measure]
    }

    type Measure {
      id: Int
      type: String
      sensor: String
      value: Float
      timestamp: Int
    }

  `);

  var getMeasures = function(args) {
      return new Promise((resolve, reject) => {
          MeasureModel.find({sensor: args.sensor}, (error, measure) => {
            resolve(measure);
          });
      });
  }

  // Root resolver
  var root = {
      measures: getMeasures
  };
  // Create an express server and a GraphQL endpoint
  var app = express();
  app.use('/graphql', express_graphql({
      schema: schema,
      rootValue: root,
      graphiql: true
  }));
  app.listen(3000, () => console.log('Express GraphQL Server Now Running On ip:3000/graphql'));

});
