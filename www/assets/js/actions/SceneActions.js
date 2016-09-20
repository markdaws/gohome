var Constants = require('../constants.js');
var Api = require('../utils/API.js');

var SceneActions = {
    create: function(scene) {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_CREATE });

            Api.sceneCreate(scene, function(err, data) {
                if (err) {
                    dispatch({ type: Constants.SCENE_CREATE_FAIL, err: err });
                    return;
                }
                dispatch({ type: Constants.SCENE_CREATE_RAW, data: data });
            });
        };
    },

    destroyClient: function() {
        return function(dispatch) {
            dispatch({ type: Constants.SCENE_DESTROY, id: "" });
            dispatch({ type: Constants.SCENE_DESTROY_RAW, id: "" });
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

    newClient: function() {
        return {
            type: Constants.SCENE_NEW_CLIENT
        };
    }
};
module.exports = SceneActions;
