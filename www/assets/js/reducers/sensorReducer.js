var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = [];

    switch(action.type) {
    case Constants.SENSOR_LOAD_ALL:
        break;

    case Constants.SENSOR_LOAD_ALL_RAW:
        newState = action.data;
        break;

    case Constants.SENSOR_LOAD_ALL_FAIL:
        //TODO: Loading error
        break;

    case Constants.SENSOR_IMPORT:
        break;
    case Constants.SENSOR_IMPORT_RAW:
        newState = [action.data].concat(newState);
        break;
    case Constants.SENSOR_IMPORT_FAIL:
        break;

    default:
        newState = state || initialState().sensors;
    }

    return newState;
};
