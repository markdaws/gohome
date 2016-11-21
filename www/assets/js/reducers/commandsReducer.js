var Constants = require('../constants.js');
var Uuid = require('uuid');

module.exports = function(state, action) {
    var newState = state;

    switch(action.type) {
    case Constants.SCENE_COMMAND_ADD:
        newState = newState.concat([action.cmd]);
        break;

    case Constants.SCENE_COMMAND_SAVE:
        break;
    case Constants.SCENE_COMMAND_SAVE_RAW:
        newState = newState.map(function(cmd) {
            if (action.cmdClientId && (cmd.clientId === action.cmdClientId)) {
                return action.cmdJson;
            }
            return cmd;
        });
        break;

    case Constants.SCENE_COMMAND_DELETE_RAW:
        newState = newState.filter(function(cmd) {
            if (action.cmdClientId && (cmd.clientId === action.cmdClientId)) {
                return false;
            }
            if (action.cmdId && (cmd.id === action.cmdId)) {
                return false;
            }
            return true;
        });
        break;
    default:
        newState = state || [];
        break;
    }
    return newState;
};
