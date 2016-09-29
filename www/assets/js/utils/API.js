var Constants = require('../constants.js');

/*
 API provides helper methods to access all of the gohome REST APIs
*/
var API = {
    // sceneLoadAll loads all of the scenes from the backing store
    sceneLoadAll: function(callback) {
        $.ajax({
            url: '/api/v1/systems/123/scenes',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr,
                });
            }
        });
    },

    // sceneCreate creates a new scene in the backing store
    sceneCreate: function(scene, callback) {
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
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (JSON.parse(xhr.responseText) || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors,
                    clientId: scene.clientId,
                });
            }
        });
    },

    // sceneUpdate updates fields of an existing scene
    sceneUpdate: function(scene, callback) {
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + scene.id,
            type: 'PUT',
            dataType: 'json',
            data: JSON.stringify(scene),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (JSON.parse(xhr.responseText) || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors,
                    sceneId: scene.id
                });
            }
        });
    },

    // sceneDestroy deletes the scene with the specified ID from the backing store
    sceneDestroy: function(id, callback) {
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + id,
            type: 'DELETE',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr
                });
            }
        });
    },

    sceneAddCommand: function(sceneId, cmd, callback) {
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + sceneId + '/commands',
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(cmd),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (JSON.parse(xhr.responseText) || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors
                });
            }
        });
    },

    sceneDeleteCommand: function() {
        //TODO: implement
    },

    // zoneLoadAll loads all of the zones from the backing store
    zoneLoadAll: function(callback) {
        $.ajax({
            url: '/api/v1/systems/123/zones',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr,
                });
            }
        });
    },

    deviceLoadAll: function(callback) {
        $.ajax({
            url: '/api/v1/systems/123/devices',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }
        });
    },

    buttonLoadAll: function(callback) {
        $.ajax({
            url: '/api/v1/systems/123/buttons',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            }.bind(this),
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }.bind(this)
        });
    }
};
module.exports = API;
