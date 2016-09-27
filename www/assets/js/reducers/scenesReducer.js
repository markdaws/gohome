var Constants = require('../constants.js');
var initialState = require('../initialState.js');

var clientId = 1;

module.exports = function(state, action) {
    var newState = Object.assign({}, state);

    switch(action.type) {
    case Constants.SCENE_LOAD_ALL:
        newState.loading = true;
        break;

    case Constants.SCENE_LOAD_ALL_FAIL:
        newState.loading = false;
        //TODO: Log fail in the UI
        break;

    case Constants.SCENE_LOAD_ALL_RAW:
        newState.loading = false;
        newState.items = action.data;
        break;

    case Constants.SCENE_NEW_CLIENT:
        newState.newSceneInfo = {
            scene: { clientId: 'scene_cid_' + clientId + '' }
        };
        newState.saveStatus = '';
        newState.saveErr = null;
        ++clientId;
        break;

    case Constants.SCENE_CREATE:
        newState.saveErr = null;
        newState.saveStatus = 'saving';
        break;

    case Constants.SCENE_CREATE_RAW:
        newState.saveStatus = 'success';
        newState.items.unshift(action.data);
        newState.newSceneInfo = null;
        break;

    case Constants.SCENE_CREATE_FAIL:
        newState.saveStatus = 'error';
        newState.saveErr = action.err;
        break;

    case Constants.SCENE_UPDATE:
        newState.saveErr = null;
        newState.saveStatus = 'saving';
        break;
    case Constants.SCENE_UPDATE_RAW:
        newState.saveStatus = 'success';
        break;
    case Constants.SCENE_UPDATE_FAIL:
        newState.saveStatus = 'error';
        newState.saveErr = action.err;
        break;

    case Constants.SCENE_DESTROY:
        break;
    case Constants.SCENE_DESTROY_RAW:
        // This is a client scene, before it was sent to the server
        if (action.id === "") {
            newState.newSceneInfo = null;
        }
        else {
            for (var i=0; i<newState.items.length; ++i) {
                if (newState.items[i].id === action.id) {
                    newState.items.splice(i, 1);
                    break;
                }
            }
        }
        break;
    case Constants.SCENE_DESTROY_FAIL:
        //TODO:
        break;

    default:
        newState = state || initialState().scenes;
    }

    return newState;
};
