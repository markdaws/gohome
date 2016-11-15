module.exports = function() {
    return {
        // TODO:
        // True if the app has loaded all the data on app load and is ready to use
        //dataLoaded: false,

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

        // An array of the zone items
        zones: [],

        // An array of all the button items
        buttons: [],

        // An array of all the sensors in the system
        sensors: [],

        // Initial load of the app
        appLoadStatus: {
            devicesLoaded: false,
            zonesLoaded: false,
            scenesLoaded: false,
            buttonsLoaded: false,
            sensorsLoaded: false
        },

        // An array of errors that should be displayed in the app
        errors: []
    };
};
