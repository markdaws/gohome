var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = state;

    switch(action.type) {
    case Constants.ERROR:
        newState = [action.data].concat(newState);
        break;

    default:
        newState = state || initialState().errors;
    }

    return newState;
};
