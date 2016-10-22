var Constants = require('../constants.js');

module.exports = function(state, action) {
    var newState = Object.assign({}, state);

    switch(action.type) {
    case Constants.SCENE_LOAD_ALL_RAW:
        newState.scenesLoaded = true;
        break;

    case Constants.ZONE_LOAD_ALL_RAW:
        newState.zonesLoaded = true;
        break;

    case Constants.DEVICE_LOAD_ALL_RAW:
        newState.devicesLoaded = true;
        break;

    case Constants.BUTTON_LOAD_ALL:
        newState.buttonsLoaded = true;
        break;
    }

    return newState;
};
