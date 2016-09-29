var Constants = require('../constants.js');

module.exports = function(state, action) {
    var newState = state;

    switch(action.type) {
    case Constants.SCENE_COMMAND_ADD:
        debugger;
        newState = [{
                    isNew: true,
                    type: action.cmdType,
                    attributes: {}
        }].concat(newState);
        break;

    default:
        newState = state || [];
        break;
    }
    return newState;
};
