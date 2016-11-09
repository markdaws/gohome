var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = state;

    switch(action.type) {
    case Constants.SENSOR_LOAD_ALL:
        break;

    case Constants.SENSOR_LOAD_ALL_RAW:
        newState = action.data;
        break;

    case Constants.SENSOR_LOAD_ALL_FAIL:
        //TODO: Loading error
        break;

    case Constants.SENSOR_UPDATE:
        break;

    case Constants.SENSOR_UPDATE_RAW:
        newState = newState.map(function(sensor) {
            if (action.data.id === sensor.id) {
                return action.data;
            }
            return sensor;
        });
        break;

    case Constants.SENSOR_UPDATE_FAIL:
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
