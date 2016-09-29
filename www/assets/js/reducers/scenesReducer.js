var Constants = require('../constants.js');
var initialState = require('../initialState.js');
var CommandsReducer = require('./commandsReducer.js');

var _clientId = 1;

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
        newState.items = [{
            clientId: 'scene_cid_' + _clientId + ''
        }].concat(newState.items);
        ++_clientId;
        break;

    case Constants.SCENE_CREATE:
        newState.saveState = Object.assign({}, newState.saveState);
        newState.saveState[action.clientId] = {
            err: null,
            status: 'saving'
        };
        break;

    case Constants.SCENE_CREATE_RAW:
        newState.saveState = Object.assign({}, newState.saveState);
        newState.saveState[action.clientId].status = 'success';

        newState.items = newState.items.map(function(scene) {
            // Replace with actual scene from the server
            if (scene.clientId === action.clientId) {
                return action.data;
            }
            return scene;
        });
        break;

    case Constants.SCENE_CREATE_FAIL:
        newState.saveState = Object.assign({}, newState.saveState);
        newState.saveState[action.clientId] = {
            status: 'error',
            err: action.err
        };
        break;

    case Constants.SCENE_UPDATE:
        newState.saveState = Object.assign({}, newState.saveState);
        newState.saveState[action.id] = {
            status: 'saving',
            err: null
        };
        break;

    case Constants.SCENE_UPDATE_RAW:
        newState.saveState = Object.assign({}, newState.saveState);
        newState.saveState[action.id].status = 'success';
        newState.items = newState.items.map(function(scene) {
            // Replace with actual scene from the server
            if (scene.id === action.id) {
                return action.sceneJson;
            }
            return scene;
        });

        break;
    case Constants.SCENE_UPDATE_FAIL:
        newState.saveState = Object.assign({}, newState.saveState);
        newState.saveState[action.id] = {
            status: 'error',
            err: action.err
        };
        break;

    case Constants.SCENE_DESTROY:
        break;
    case Constants.SCENE_DESTROY_RAW:
        // This is a client scene, before it was sent to the server
        for (var i=0; i<newState.items.length; ++i) {
            var found = false;
            if (action.clientId) {
                found = newState.items[i].clientId === action.clientId;
            } else {
                found = newState.items[i].id === action.id;
            }

            if (found) {
                newState.items = newState.items.slice();
                newState.items.splice(i, 1);
                break;
            }
        }
        break;
    case Constants.SCENE_DESTROY_FAIL:
        //TODO:
        break;


    case Constants.SCENE_COMMAND_ADD:
        
        /*
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
        }*/
        break;
    case Constants.SCENE_COMMAND_ADD_RAW:
        break;
    case Constants.SCENE_COMMAND_ADD_FAIL:
        break;
    default:
        newState = state || initialState().scenes;
    }

    return newState;
};
