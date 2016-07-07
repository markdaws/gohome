var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var Constants = require('../constants.js');
var Api = require('../utils/API.js');

var SceneActions = {
    create: function(scene) {
        AppDispatcher.dispatch({
            actionType: Constants.SCENE_CREATE,
            scene: scene,
        });
        Api.sceneCreate(scene);
    },

    destroy: function(id) {
        AppDispatcher.dispatch({
            actionType: Constants.SCENE_DESTROY,
            id: id,
        });
        Api.sceneDestroy(id);
    },

    loadAll: function() {
        return function(dispatch) {
            dispatch({
                type: Constants.SCENE_LOAD_ALL
            });

            Api.sceneLoadAll(function(error, data) {
                //error.err / error.status / error.xhr
                if (error) {
                    dispatch({ type: Constants.SCENE_LOAD_ALL_FAIL, err: error });
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
