var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = Object.assign({}, state);

    switch(action.type) {
    case Constants.DEVICE_LOAD_ALL:
        break;

    case Constants.DEVICE_LOAD_ALL_RAW:
        newState.devices = action.data;
        break;

    case Constants.DEVICE_LOAD_ALL_FAIL:
        //TODO: Loading error
        break;

    default:
        newState = state || initialState().system;
    }

    return newState;
};
