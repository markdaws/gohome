var Redux = require('redux');
var thunk = require('redux-thunk').default;
var initialState = require('./initialState.js');
var buttonReducer = require('./reducers/buttonReducer.js');
var systemReducer = require('./reducers/systemReducer.js');
var scenesReducer = require('./reducers/scenesReducer.js');
var zonesReducer = require('./reducers/zonesReducer.js');

var rootReducer = Redux.combineReducers({
    system: systemReducer,
    scenes: scenesReducer,
    zones: zonesReducer,
    buttons: buttonReducer
});

module.exports = Redux.applyMiddleware(thunk)(Redux.createStore)(rootReducer, initialState());
