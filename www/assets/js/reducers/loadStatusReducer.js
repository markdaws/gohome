var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = state;

    switch(action.type) {
    case Constants.SCENE_LOAD_ALL_RAW:
        newState = Object.assign({}, state);
        newState.scenesLoaded = true;
        break;

    case Constants.DEVICE_LOAD_ALL_RAW:
        newState = Object.assign({}, state);
        newState.devicesLoaded = true;
        break;

    case Constants.AUTOMATION_LOAD_ALL_RAW:
        newState = Object.assign({}, state);
        newState.automationLoaded = true;
        break;

    default:
        newState = state || initialState().appLoadStatus;
    }

    return newState;
};
