var Constants = require('../constants.js');

module.exports = function(state, action) {
    var newState;

    switch(action.type) {
    case Constants.SCENE_COMMAND_ADD:
        debugger;
        var scenes = newState.items;
        for (var i=0;i<scenes.length; ++i) {
            if (scenes[i].id === action.sceneId) {
                scenes[i].commands.push({
                    isNew: true,
                    type: action.cmdType,
                    attributes: {}
                });
                break;
            }
        }
        break;

    case Constants.SCENE_COMMAND_ADD_RAW:
        break;

    case Constants.SCENE_COMMAND_ADD_FAIL:
        break;

    default:
        newState = state || [];
        break;
    }
    return newState;
};
