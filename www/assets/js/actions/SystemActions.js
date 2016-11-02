var Constants = require('../constants.js');
var Api = require('../utils/API.js');

var SystemActions = {
    deviceNew: function() {
        return function(dispatch) {
            dispatch({ type: Constants.DEVICE_NEW_CLIENT });
        };
    },

    deviceDelete: function(id, clientId) {
        return function(dispatch) {
            dispatch({ type: Constants.DEVICE_DESTROY });
            if (!id) {
                dispatch({ type: Constants.DEVICE_DESTROY_RAW, clientId: clientId });
                return;
            }

            Api.deviceDestroy(id, function(err, data) {
                if (err) {
                    dispatch({ type: Constants.DEVICE_DESTROY_FAIL, id: id, clientId: clientId, err: err });
                    return;
                }
                dispatch({ type: Constants.DEVICE_DESTROY_RAW, id: id, clientId: clientId, data: data });
            });
        };
    },

    createdDevice: function(clientId, deviceJson, append) {
        return function(dispatch) {
            dispatch({ type: Constants.DEVICE_CREATE_RAW, data: deviceJson, clientId: clientId });
        };
    },

    updatedDevice: function(deviceJson) {
        return function(dispatch) {
            dispatch({ type: Constants.DEVICE_UPDATE_RAW, data: deviceJson });
        };
    },

    importedDevice: function(deviceJson) {
        return function(dispatch) {
            dispatch({ type: Constants.DEVICE_IMPORT_RAW, data: deviceJson });
        };
    },

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

            Api.buttonLoadAll(function(err, data) {
                if (err) {
                    dispatch({ type: Constants.BUTTON_LOAD_ALL_FAIL, err: err });
                    return;
                }
                dispatch({ type: Constants.BUTTON_LOAD_ALL_RAW, data: data });
            });
        };
    }
};
module.exports = SystemActions;
