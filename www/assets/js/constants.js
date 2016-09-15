var keyMirror = require('keyMirror');

/*
 The pattern here is that you have a actionType e.g. SCENE_DESTROY, responses from the
 server are then either <ACTION_TYPE>_RAW which indicates success and the payload will
 contain any raw data returned from the server, or <ACTION_TYPE>_FAIL which indicates
 failure
 */

module.exports = keyMirror({

    // Load all of the scenes from the server
    SCENE_LOAD_ALL: null,
    SCENE_LOAD_ALL_RAW: null,
    SCENE_LOAD_ALL_FAIL: null,

    // Create a new scene on the client
    SCENE_NEW_CLIENT: null,

    // Creates a new scene on the server
    SCENE_CREATE: null,
    SCENE_CREATE_RAW: null,
    SCENE_CREATE_FAIL: null,

    // Update the attribute of a scene
    SCENE_UPDATE: null,
    SCENE_UPDATE_RAW: null,
    SCENE_UPDATE_FAIL: null,

    // Delete a scene
    SCENE_DESTROY: null,
    SCENE_DESTROY_RAW: null,
    SCENE_DESTROY_FAIL: null,

    // Load all of the zones from the server
    ZONE_LOAD_ALL: null,
    ZONE_LOAD_ALL_RAW: null,
    ZONE_LOAD_ALL_FAIL: null,
});
