var Redux = require('redux');
var thunk = require('redux-thunk').default;
var initialState = require('./initialState.js');
var scenesReducer = require('./reducers/scenesReducer.js');
var zonesReducer = require('./reducers/zonesReducer.js');

var rootReducer = Redux.combineReducers({
    scenes: scenesReducer,
    zones: zonesReducer
});

module.exports = Redux.applyMiddleware(thunk)(Redux.createStore)(rootReducer, initialState());
