var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var Constants = require('../constants/constants.js');
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
        AppDispatcher.dispatch({
            actionType: Constants.SCENE_LOAD_ALL,
        });
        Api.sceneLoadAll();
    },
    
    newClient: function() {
        AppDispatcher.dispatch({
            actionType: Constants.SCENE_NEW_CLIENT,
        });
    }
};
module.exports = SceneActions;
