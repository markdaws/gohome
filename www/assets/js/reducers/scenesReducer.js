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
            scene: { clientId: 'scene_cid_' + clientId + '' },
            saveErr: null,
            saveStatus: null
        };
        ++clientId;
        break;

    case Constants.SCENE_CREATE:
        newState.newSceneInfo.saveStatus = 'saving';
        break;

    case Constants.SCENE_CREATE_RAW:
        newState.newSceneInfo.saveStatus = 'success';
        // TODO: Refetch all scenes
        break;

    case Constants.SCENE_CREATE_FAIL:
        newState.newSceneInfo.saveStatus = 'error';
        newState.newSceneInfo.saveErr = action.err;
        break;

    case Constants.SCENE_DESTROY:
        break;
    case Constants.SCENE_DESTROY_RAW:
        
        for (var i=0; i<newState.items.length; ++i) {
            if (newState.items[i].id === action.id) {
                newState.items.splice(i, 1);
                break;
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
