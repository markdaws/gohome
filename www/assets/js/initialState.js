module.exports = function() {
    return {
        scenes: {
            // true if currently loading the scene list
            loading: false,

            //TODO: loadingErr

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
        zones: {
            // true if currently loading the zone list
            loading: false,

            // contains the error object if the zones failed to load
            loadingErr: null,

            // array of zone objects
            items: []
        }
    };
};
