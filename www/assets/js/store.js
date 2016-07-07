var Redux = require('redux');
var thunk = require('redux-thunk').default;
var initialState = require('./initialState.js');
var scenesReducer = require('./reducers/scenesReducer.js');

var rootReducer = Redux.combineReducers({
    scenes: scenesReducer
});

module.exports = Redux.applyMiddleware(thunk)(Redux.createStore)(rootReducer, initialState());
