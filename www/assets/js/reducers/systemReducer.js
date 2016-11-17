var Constants = require('../constants.js');
var initialState = require('../initialState.js');
var uuid = require('uuid');

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

    case Constants.DEVICE_NEW_CLIENT:
        newState.devices = [{
            id: uuid.v4()
        }].concat(newState.devices);
        break;

    case Constants.DEVICE_CREATE:
        break;
    case Constants.DEVICE_CREATE_RAW:
        newState.devices = newState.devices.map(function(device) {
            if (device.id === action.data.id) {
                return action.data;
            }
            return device;
        });

    case Constants.DEVICE_CREATE_FAIL:
        break;

    case Constants.DEVICE_UPDATE:
        break;

    case Constants.DEVICE_UPDATE_RAW:
        newState.devices = newState.devices.map(function(device) {
            if (device.id === action.data.id) {
                return action.data;
            }
            return device;
        });
        break;

    case Constants.DEVICE_UPDATE_FAIL:
        break;

    case Constants.DEVICE_IMPORT:
        break;
    case Constants.DEVICE_IMPORT_RAW:
        newState.devices = [action.data].concat(newState.devices);
        break;
    case Constants.DEVICE_IMPORT_FAIL:
        break;
        
    case Constants.DEVICE_DESTROY:
        break;

    case Constants.DEVICE_DESTROY_RAW:
        // This is a client device, before it was sent to the server
        for (var i=0; i<newState.devices.length; ++i) {
            var found = newState.devices[i].id === action.id;

            if (found) {
                newState.devices = newState.devices.slice();
                newState.devices.splice(i, 1);
                break;
            }
        }
        break;
        
    case Constants.DEVICE_DESTROY_FAIL:
        //TODO:
        break;

    default:
        newState = state || initialState().system;
    }

    return newState;
};
