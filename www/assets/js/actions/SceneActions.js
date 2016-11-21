var Constants = require('../constants.js');
var Api = require('../utils/API.js');

var SceneActions = {

    loadAll: function() {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_LOAD_ALL });

            Api.sceneLoadAll(function(err, data) {
                if (err) {
                    dispatch({ type: Constants.SCENE_LOAD_ALL_FAIL, err: err });
                    return;
                }

                dispatch({ type: Constants.SCENE_LOAD_ALL_RAW, data: data });
            });
        };
    },

    created: function(sceneJson, clientId) {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_CREATE_RAW, data: sceneJson, id: sceneJson.id, clientId: clientId });
        };
    },

    update: function(sceneJson) {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_UPDATE, id: sceneJson.id });

            Api.sceneUpdate(sceneJson, function(err, data) {
                if (err) {
                    dispatch({ type: Constants.SCENE_UPDATE_FAIL, err: err, id: sceneJson.id });
                    return;
                }
                dispatch({ type: Constants.SCENE_UPDATE_RAW, data: data, id: sceneJson.id, sceneJson: sceneJson });
            });
        };
    },

    destroyClient: function(clientId) {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_DESTROY, clientId: clientId });
            dispatch({ type: Constants.SCENE_DESTROY_RAW, clientId: clientId });
        };
    },

    destroy: function(id) {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_DESTROY, id: id });

            Api.sceneDestroy(id, function(err, data) {
                if (err) {
                    dispatch({ type: Constants.SCENE_DESTROY_FAIL, err: err, id: id });
                    return;
                }
                dispatch({ type: Constants.SCENE_DESTROY_RAW, data: data, id: id });
            });
        };
    },

    addCommand: function(sceneId, cmd) {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_COMMAND_ADD, sceneId: sceneId, cmd: cmd });
        };
    },

    savedCommand: function(sceneId, cmdClientId, cmdJson) {
        return function(dispatch) {
            dispatch({
                type: Constants.SCENE_COMMAND_SAVE_RAW,
                sceneId: sceneId,
                cmdClientId: cmdClientId,
                cmdJson: cmdJson
            });
        };
    },

    deleteCommand: function(sceneId, cmdId, cmdClientId) {
        return function(dispatch) {
            dispatch({
                type: Constants.SCENE_COMMAND_DELETE,
                cmdId: cmdId,
                cmdClientId: cmdClientId,
                sceneId: sceneId });

            if (cmdClientId) {
                // Client only, not saved on the server
                dispatch({
                    type: Constants.SCENE_COMMAND_DELETE_RAW,
                    sceneId: sceneId,
                    cmdClientId: cmdClientId });
                return;
            }

            Api.sceneDeleteCommand(sceneId, cmdId, function(err, data) {
                if (err) {
                    dispatch({
                        type: Constants.SCENE_COMMAND_DELETE_FAIL,
                        sceneId: sceneId,
                        cmdId: cmdId,
                        err: err });
                    return;
                }
                dispatch({
                    type: Constants.SCENE_COMMAND_DELETE_RAW,
                    sceneId: sceneId,
                    cmdId: cmdId });
            });
        };
    },

    newClient: function() {
        return {
            type: Constants.SCENE_NEW_CLIENT
        };
    }
};
module.exports = SceneActions;
