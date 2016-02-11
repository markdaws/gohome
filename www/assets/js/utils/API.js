var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var Constants = require('../constants/constants.js');

/*
 API provides helper methods to access all of the gohome REST APIs
*/
var API = {
    // sceneLoadAll loads all of the scenes from the backing store
    sceneLoadAll: function() {
        $.ajax({
            url: '/api/v1/systems/123/scenes',
            dataType: 'json',
            cache: false,
            success: function(data) {
                AppDispatcher.dispatch({
                    actionType: Constants.SCENE_LOAD_ALL_RAW,
                    raw: data,
                });
            },
            error: function(xhr, status, err) {
                AppDispatcher.dispatch({
                    actionType: Constants.SCENE_LOAD_ALL_FAIL,
                    err: err,
                    xhr: xhr,
                });
            }
        });
    },

    // sceneCreate creates a new scene in the backing store
    sceneCreate: function(scene) {
        // Note: new scenes don't have an ID yet, since the server has to assign that, but
        // they do have a clientId which is a unique ID created on the client so they can
        // still be distinguished from one another
        $.ajax({
            url: '/api/v1/systems/123/scenes',
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(scene),
            cache: false,
            success: function(data) {
                AppDispatcher.dispatch({
                    actionType: Constants.SCENE_CREATE_RAW,
                    raw: data,
                    clientId: scene.clientId,
                });
            },
            error: function(xhr, status, err) {
                var errors = (JSON.parse(xhr.responseText) || {}).errors;
                AppDispatcher.dispatch({
                    actionType: Constants.SCENE_CREATE_FAIL,
                    err: err,
                    xhr: xhr,
                    validationErrors: errors,
                    clientId: scene.clientId,
                });
            }
        });
    },

    // sceneDestroy deletes the scene with the specified ID from the backing store
    sceneDestroy: function(id) {
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + id,
            type: 'DELETE',
            cache: false,
            success: function(data) {
                AppDispatcher.dispatch({
                    actionType: Constants.SCENE_DESTROY_RAW,
                    id: id,
                });
            },
            error: function(xhr, status, err) {
                AppDispatcher.dispatch({
                    actionType: Constants.SCENE_DESTROY_FAIL,
                    id: id,
                    err: err,
                    xhr: xhr,
                });
            }
        });
    },

    // zoneLoadAll loads all of the zones from the backing store
    zoneLoadAll: function() {
        $.ajax({
            url: '/api/v1/systems/123/zones',
            dataType: 'json',
            cache: false,
            success: function(data) {
                AppDispatcher.dispatch({
                    actionType: Constants.ZONE_LOAD_ALL_RAW,
                    raw: data,
                });
            },
            error: function(xhr, status, err) {
                AppDispatcher.dispatch({
                    actionType: Constants.ZONE_LOAD_ALL_FAIL,
                    err: err,
                    xhr: xhr,
                });
            }
        });
    }
};
module.exports = API;
