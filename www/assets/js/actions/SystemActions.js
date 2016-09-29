var Constants = require('../constants.js');
var Api = require('../utils/API.js');

var SystemActions = {
    loadAllDevices: function() {
        return function(dispatch) {
            dispatch({ type: Constants.DEVICE_LOAD_ALL });

            Api.deviceLoadAll(function(err, data) {
                if (err) {
                    dispatch({ type: Constants.DEVICE_LOAD_ALL_FAIL, err: err });
                    return;
                }
                dispatch({ type: Constants.DEVICE_LOAD_ALL_RAW, data: data });
            });
        };
    },

    loadAllButtons: function() {
        return function(dispatch) {
            dispatch({ type: Constants.BUTTON_LOAD_ALL });

            Api
        };
    }
};
module.exports = SystemActions;
