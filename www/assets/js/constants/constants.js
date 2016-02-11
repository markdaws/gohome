var keyMirror = require('keyMirror');

// The pattern here is that you have a actionType e.g. SCENE_CREATE, responses from the
// server are then either <ACTION_TYPE>_RAW which indicates success and the payload will
// contain any raw data returned from the server, or <ACTION_TYPE>_FAIL which indicates
// failure
module.exports = keyMirror({
    SCENE_LOAD_ALL: null,
    SCENE_LOAD_ALL_RAW: null,
    SCENE_LOAD_ALL_FAIL: null,

    SCENE_CREATE: null,
    SCENE_CREATE_RAW: null,
    SCENE_CREATE_FAIL: null,
    SCENE_NEW_CLIENT: null,
    
    SCENE_DESTROY: null,
    SCENE_DESTROY_RAW: null,
    SCENE_DESTROY_FAIL: null,
    
    //TODO:
    SCENE_UPDATE: null,
    SCENE_CREATE_COMMAND: null,
    SCENE_DESTROY_COMMAND: null,

    ZONE_LOAD_ALL: null,
    ZONE_LOAD_ALL_RAW: null,
    ZONE_LOAD_ALL_FAIL: null,
});
