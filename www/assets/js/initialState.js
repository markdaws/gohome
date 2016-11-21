module.exports = function() {
    return {
        system: {
            // Array of devices in the system
            devices: []
        },

        scenes: {
            // array of scene objects
            items: [],

            // TODO: Rethink
            // Save state of the different scenes, will be keyed by id, or  client id if no id
            saveState: { }
        },

        // Initial load of the app
        appLoadStatus: {
            devicesLoaded: false,
            scenesLoaded: false,
        },

        // An array of errors that should be displayed in the app
        errors: []
    };
};
