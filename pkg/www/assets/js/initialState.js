module.exports = function() {
    return {
        system: {
            // Array of devices in the system
            devices: []
        },

        scenes: [],

        // Initial load of the app
        appLoadStatus: {
            devicesLoaded: false,
            scenesLoaded: false,
            automationLoaded: false
        },

        // Array of all the automation scripts loaded in the system
        automations: [],
    };
};
