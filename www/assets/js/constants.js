var keyMirror = require('keyMirror');

/*
 The pattern here is that you have a actionType e.g. SCENE_DESTROY, responses from the
 server are then either <ACTION_TYPE>_RAW which indicates success and the payload will
 contain any raw data returned from the server, or <ACTION_TYPE>_FAIL which indicates
 failure. If the errors are handled localy inside a component, there may not be a corresponding
 _FAIL message
 */

module.exports = keyMirror({

    // Load all of the automation items from the server
    AUTOMATION_LOAD_ALL: null,
    AUTOMATION_LOAD_ALL_RAW: null,
    AUTOMATION_LOAD_ALL_FAIL: null,

    // Load all of the devices from the server
    DEVICE_LOAD_ALL: null,
    DEVICE_LOAD_ALL_RAW: null,
    DEVICE_LOAD_ALL_FAIL: null,

    // Adds a new device on the client
    DEVICE_NEW_CLIENT: null,

    // Creates a new device on the server
    DEVICE_CREATE: null,
    DEVICE_CREATE_RAW: null,
    DEVICE_CREATE_FAIL: null,

    // Device fields have been updated
    DEVICE_UPDATE: null,
    DEVICE_UPDATE_RAW: null,
    DEVICE_UPDATE_FAIL: null,

    // When we are importing a device
    DEVICE_IMPORT: null,
    DEVICE_IMPORT_RAW: null,
    DEVICE_IMPORT_FAIL: null,

    // Deletes a device
    DEVICE_DESTROY: null,
    DEVICE_DESTROY_RAW: null,
    DEVICE_DESTROY_FAIL: null,

    // A new global error has been fired
    ERROR: null,

    // Clears all of the global errors
    ERROR_CLEAR: null,

    // Updated a feature
    FEATURE_UPDATE: null,
    FEATURE_UPDATE_RAW: null,
    FEATURE_UPDATE_FAIL: null,

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

    // Add a command to a scene, not saved to the server, just on the client
    SCENE_COMMAND_ADD: null,

    // Saves a command associated to a scene on the server
    SCENE_COMMAND_SAVE: null,
    SCENE_COMMAND_SAVE_RAW: null,

    // Remove a command from a scene
    SCENE_COMMAND_DELETE: null,
    SCENE_COMMAND_DELETE_RAW: null,
    SCENE_COMMAND_DELETE_FAIL: null,
});
