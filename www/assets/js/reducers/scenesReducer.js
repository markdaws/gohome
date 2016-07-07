var Api = require('../utils/API.js');
var Constants = require('../Constants.js');
var initialState = require('../initialState.js');

var clientId = 1;

module.exports = function(state, action) {
    var newState = Object.assign({}, state);

    switch(action.type) {
    case Constants.SCENE_LOAD_ALL:
        newState.loading = true;
        break;

    case Constants.SCENE_LOAD_ALL_FAIL:
        console.log('scene load all fail');
        console.log(action);

        newState.loading = false;
        //TODO: Log fail in the UI
        break;

    case Constants.SCENE_LOAD_ALL_RAW:
        console.log('scene load all raw');
        console.log(action);

        newState.loading = false;
        newState.items = action.data;
        break;

    case Constants.SCENE_NEW_CLIENT:
        console.log('scene new client');

        newState.items.push({ clientId: 'scene_cid_' + clientId + '' });
        ++clientId;
        break;

    default:
        newState = state || initialState().scenes;
    }

    return newState;
};
