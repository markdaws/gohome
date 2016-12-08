var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = state;

    switch(action.type) {
    case Constants.AUTOMATION_LOAD_ALL:
        break;

    case Constants.AUTOMATION_LOAD_ALL_RAW:
        action.data.sort(function(a, b) {
            return a.name.localeCompare(b.name);
        });
        newState = action.data;
        break;

    case Constants.AUTOMATION_LOAD_ALL_FAIL:
        break;
    default:
        newState = newState ||initialState().automations;
    }

    return newState;
};
