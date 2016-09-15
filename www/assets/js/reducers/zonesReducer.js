var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = Object.assign({}, state);

    switch(action.type) {
    case Constants.ZONE_LOAD_ALL:
        newState.loading = true;
        break;

    case Constants.ZONE_LOAD_ALL_FAIL:
        newState.loading = false;
        newState.loadingErr = action.err;
        break;

    case Constants.ZONE_LOAD_ALL_RAW:
        newState.loading = false;
        newState.items = action.data;
        break;

    default:
        newState = state || initialState();
    }

    return newState;
};
