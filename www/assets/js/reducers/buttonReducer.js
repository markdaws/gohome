var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = [];

    switch(action.type) {
    case Constants.BUTTON_LOAD_ALL:
        break;

    case Constants.BUTTON_LOAD_ALL_RAW:
        newState = action.data;
        break;

    case Constants.BUTTON_LOAD_ALL_FAIL:
        //TODO: Loading error
        break;

    default:
        newState = state || initialState().buttons;
    }

    return newState;
};
