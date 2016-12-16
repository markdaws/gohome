var Constants = require('../constants.js');
var initialState = require('../initialState.js');
var CommandsReducer = require('./commandsReducer.js');
var  uuid = require('uuid');

module.exports = function(state, action) {
    var newState = state;
    var i;

    switch(action.type) {
    case Constants.SCENE_LOAD_ALL:
        break;

    case Constants.SCENE_LOAD_ALL_FAIL:
        //TODO: Log fail in the UI
        break;

    case Constants.SCENE_LOAD_ALL_RAW:
        action.data.sort(function(a, b) {
            return a.name.localeCompare(b.name);
        });
        newState = action.data;
        break;

    case Constants.SCENE_NEW_CLIENT:
        newState = [{
            clientId: uuid.v4()
        }].concat(newState);
        break;

    case Constants.SCENE_CREATE_RAW:
        newState = newState.map(function(scene) {
            // Replace with actual scene from the server
            if (scene.clientId === action.clientId) {
                delete action.data.clientId;
                return action.data;
            }
            return scene;
        });
        break;

    case Constants.SCENE_UPDATE_RAW:
        newState = newState.map(function(scene) {
            // Replace with actual scene from the server
            if (scene.id === action.id) {
                return action.data;
            }
            return scene;
        });
        break;

    case Constants.SCENE_DESTROY:
        break;

    case Constants.SCENE_DESTROY_RAW:
        for (i=0; i<newState.length; ++i) {
            var found = false;
            if (action.clientId && (action.clientId === newState[i].clientId)) {
                found = true;
            }
            if (action.id && (action.id === newState[i].id)) {
                found = true;
            }
            if (found) {
                newState = newState.slice();
                newState.splice(i, 1);
                break;
            }
        }
        break;

    case Constants.SCENE_DESTROY_FAIL:
        //TODO:
        break;

    case Constants.SCENE_COMMAND_ADD:
    case Constants.SCENE_COMMAND_SAVE:
    case Constants.SCENE_COMMAND_SAVE_RAW:
    case Constants.SCENE_COMMAND_SAVE_FAIL:
    case Constants.SCENE_COMMAND_DELETE_RAW:
        var scenes = newState;
        for (i=0;i<scenes.length; ++i) {
            if (scenes[i].id === action.sceneId) {
                newState = newState.slice();
                newState[i].commands = CommandsReducer(scenes[i].commands || [], action);
                break;
            }
        }
        break;

    default:
        newState = state || initialState().scenes;
    }

    return newState;
};
