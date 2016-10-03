var Constants = require('../constants.js');

module.exports = function(state, action) {
    var newState = state;

    switch(action.type) {
    case Constants.SCENE_COMMAND_ADD:
        newState = newState.concat([{
                    isNew: true,
                    type: action.cmdType,
                    attributes: {}
        }]);
        break;

    case Constants.SCENE_COMMAND_SAVE:
        break;
    case Constants.SCENE_COMMAND_SAVE_RAW:
        newState = newState.slice();
        newState[action.cmdIndex] = Object.assign({}, action.data, { isNew: false });
        break;

    case Constants.SCENE_COMMAND_DELETE_RAW:
        newState = newState.slice();
        newState.splice(action.cmdIndex, 1);
        break;
    default:
        newState = state || [];
        break;
    }
    return newState;
};
