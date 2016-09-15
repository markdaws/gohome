module.exports = function() {
    return {
        scenes: {
            // true if currently loading the scene list
            loading: false,

            // if the user is creating a new scene on the client, this value
            // will be populated as on object with fields scene/saveErr
            newSceneInfo: null,

            // array of scene objects
            items: []
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
