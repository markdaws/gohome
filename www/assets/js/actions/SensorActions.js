var Constants = require('../constants.js');
var Api = require('../utils/API.js');

var SensorActions = {
    loadAll: function() {
        return function(dispatch) {
            dispatch({
                type: Constants.SENSOR_LOAD_ALL
            });

            Api.sensorLoadAll(function(err, data) {
                if (err) {
                    dispatch({ type: Constants.SENSOR_LOAD_ALL_FAIL, err: err });
                    return;
                }

                dispatch({ type: Constants.SENSOR_LOAD_ALL_RAW, data: data });
            });
        };
    },

    importedSensor: function(sensorJson) {
        return function(dispatch) {
            dispatch({ type: Constants.SENSOR_IMPORT_RAW, data: sensorJson });
        };
    },

};
module.exports = SensorActions;
