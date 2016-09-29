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
            // Save status, '""|"saving"|"success"|"error"', can be saving of a new scene
            // or saving of an update to an existing scene
            saveStatus: null,

            // Detailed object with more description on the save error
            saveErr: null,

            // array of scene objects
            items: [],

            // Save state of the different scenes, will be keyed by id, or  client id if no id
            saveState: { }
        },

        // An array of the zone items
        zones: [],

        // An array of all the button items
        buttons: []
    };
};
