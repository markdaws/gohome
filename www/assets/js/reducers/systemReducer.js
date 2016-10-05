var Constants = require('../constants.js');
var initialState = require('../initialState.js');

var _clientId = 1;
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
            clientId: 'device_cid_' + _clientId
        }].concat(newState.devices);
        ++_clientId;
        break;

    case Constants.DEVICE_CREATE:
        break;
    case Constants.DEVICE_CREATE_RAW:
        break;
    case Constants.DEVICE_CREATE_FAIL:
        break;

    case Constants.DEVICE_DESTROY:
        break;

    case Constants.DEVICE_DESTROY_RAW:
        debugger;
        // This is a client device, before it was sent to the server
        for (var i=0; i<newState.devices.length; ++i) {
            var found = false;
            if (action.id) {
                found = newState.devices[i].id === action.id;
            } else {
                found = newState.devices[i].clientId === action.clientId;
            }

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
