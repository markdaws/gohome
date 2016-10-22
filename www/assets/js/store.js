var Redux = require('redux');
var thunk = require('redux-thunk').default;
var initialState = require('./initialState.js');
var buttonReducer = require('./reducers/buttonReducer.js');
var systemReducer = require('./reducers/systemReducer.js');
var scenesReducer = require('./reducers/scenesReducer.js');
var zonesReducer = require('./reducers/zonesReducer.js');
var loadStatusReducer = require('./reducers/loadStatusReducer.js');


var rootReducer = Redux.combineReducers({
    system: systemReducer,
    scenes: scenesReducer,
    zones: zonesReducer,
    buttons: buttonReducer,
    appLoadStatus: loadStatusReducer
});

module.exports = Redux.applyMiddleware(thunk)(Redux.createStore)(rootReducer, initialState());
