var Constants = require('../constants.js');
var Api = require('../utils/API.js');

var ZoneActions = {
    loadAll: function() {
        return function(dispatch) {
            dispatch({
                type: Constants.ZONE_LOAD_ALL
            });

            Api.zoneLoadAll(function(err, data) {
                if (err) {
                    dispatch({ type: Constants.ZONE_LOAD_ALL_FAIL, err: err });
                    return;
                }

                dispatch({ type: Constants.ZONE_LOAD_ALL_RAW, data: data });
            });
        };
    },

    importedZone: function(zoneJson) {
        return function(dispatch) {
            dispatch({ type: Constants.ZONE_IMPORT_RAW, data: zoneJson });
        };
    },

    updated: function(zoneJson) {
        return function(dispatch) {
            dispatch({ type: Constants.ZONE_UPDATE_RAW, data: zoneJson });
        };
    }
};
module.exports = ZoneActions;
