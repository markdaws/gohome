var Redux = require('redux');
var thunk = require('redux-thunk').default;
var initialState = require('./initialState.js');
var systemReducer = require('./reducers/systemReducer.js');
var scenesReducer = require('./reducers/scenesReducer.js');
var loadStatusReducer = require('./reducers/loadStatusReducer.js');
var errorReducer = require('./reducers/errorReducer.js');
var automationReducer = require('./reducers/automationReducer.js');

var rootReducer = Redux.combineReducers({
    system: systemReducer,
    scenes: scenesReducer,
    appLoadStatus: loadStatusReducer,
    errors: errorReducer,
    automations: automationReducer,
});

module.exports = Redux.applyMiddleware(thunk)(Redux.createStore)(rootReducer, initialState());
