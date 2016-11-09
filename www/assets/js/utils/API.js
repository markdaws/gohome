var Constants = require('../constants.js');

var BASE = '//' + window.location.hostname + ':5000';

/*
 API provides helper methods to access all of the gohome REST APIs
*/
var API = {

    //TODO: Break this in to separate files for different objects

    deviceLoadAll: function(callback) {
        $.ajax({
            url: BASE + '/api/v1/devices',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }
        });
    },

    deviceCreate: function(deviceJson, callback) {
        $.ajax({
            url: BASE + '/api/v1/devices',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(deviceJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors
                });
            }
        });
    },

    deviceUpdate: function(deviceJson, callback) {
        $.ajax({
            url: BASE + '/api/v1/devices/' + deviceJson.id,
            type: 'PUT',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(deviceJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors
                });
            }
        });        
    },

    // Deletes a device on the server
    deviceDestroy: function(id, callback) {
        $.ajax({
            url: BASE + '/api/v1/devices/' + id,
            type: 'DELETE',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr
                });
            }
        });
    },

    // sceneActivate actives the specified scene
    sceneActivate: function(sceneId, callback) {
        $.ajax({
            url: BASE + '/api/v1/scenes/active',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify({ id: sceneId }),
            success: function(data) {
                callback(null, data);
            }.bind(this),
            error: function(xhr, status, err) {
                callback({ err: err });
            }.bind(this)
        });
    },

    // sceneLoadAll loads all of the scenes from the backing store
    sceneLoadAll: function(callback) {
        $.ajax({
            url: BASE + '/api/v1/scenes',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr,
                });
            }
        });
    },

    // sceneCreate creates a new scene in the backing store
    sceneCreate: function(scene, callback) {
        // Note: new scenes don't have an ID yet, since the server has to assign that, but
        // they do have a clientId which is a unique ID created on the client so they can
        // still be distinguished from one another
        $.ajax({
            url: BASE + '/api/v1/scenes',
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(scene),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors,
                    clientId: scene.clientId
                });
            }
        });
    },

    // sceneUpdate updates fields of an existing scene
    sceneUpdate: function(scene, callback) {
        $.ajax({
            url: BASE + '/api/v1/scenes/' + scene.id,
            type: 'PUT',
            dataType: 'json',
            data: JSON.stringify(scene),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors,
                    id: scene.id
                });
            }
        });
    },

    // sceneDestroy deletes the scene with the specified ID from the backing store
    sceneDestroy: function(id, callback) {
        $.ajax({
            url: BASE + '/api/v1/scenes/' + id,
            type: 'DELETE',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr
                });
            }
        });
    },

    sceneSaveCommand: function(sceneId, cmd, callback) {
        $.ajax({
            url: BASE + '/api/v1/scenes/' + sceneId + '/commands',
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(cmd),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors
                });
            }
        });
    },

    sceneDeleteCommand: function(sceneId, cmdIndex, callback) {
        $.ajax({
            url: BASE + '/api/v1/scenes/' + sceneId + '/commands/' + cmdIndex,
            type: 'DELETE',
            dataType: 'json',
            data: {},
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }
        });
    },

    // monitorSubscribe requests to subscribe to sensor and zone changes
    monitorSubscribe: function(group, callback) {
        $.ajax({
            url: BASE + '/api/v1/monitor/groups',
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(group),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }
        });
    },

    // monitorSubscribeRenew renews the subscription to the monitor group, this extends
    // the timeout for the group before the server stops sending updates
    monitorSubscribeRenew: function(monitorId, callback) {
        $.ajax({
            url: BASE + '/api/v1/monitor/groups/' + monitorId,
            type: 'PUT',
            dataType: 'json',
            data: null,
            cache: false,
            success: function(data) {
                callback && callback(null, data);
            },
            error: function(xhr, status, err) {
                callback && callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }
        });
    },

    // monitorUnsubscribe unsubscribe the specified monitor id so the client will
    // no longer receive updates when values associated with it change
    monitorUnsubscribe: function(monitorId, callback) {
        $.ajax({
            url: BASE + '/api/v1/monitor/groups/' + monitorId,
            type: 'DELETE',
            dataType: 'json',
            data: null,
            cache: false,
            success: function(data) {
                callback && callback(null, data);
            },
            error: function(xhr, status, err) {
                callback && callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }
        });
    },
    
    // sensorLoadAll loads all of the sensors from the backing store
    sensorLoadAll: function(callback) {
        $.ajax({
            url: BASE + '/api/v1/sensors',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr,
                });
            }
        });
    },

    // sensorCreate creates a new sensor on the server
    sensorCreate: function(sensorJson, callback) {
        $.ajax({
            url: BASE + '/api/v1/sensors',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(sensorJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    status: status,
                    validationErrors: errors
                });
            }
        });
    },

    // sensorUpdate updates a sensor on the server with the new values
    sensorUpdate: function(sensorJson, callback) {
        $.ajax({
            url: BASE + '/api/v1/sensors/' + sensorJson.id,
            type: 'PUT',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(sensorJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors
                });
            }
        });        
    },

    // zoneLoadAll loads all of the zones from the backing store
    zoneLoadAll: function(callback) {
        $.ajax({
            url: BASE + '/api/v1/zones',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    status: status,
                    xhr: xhr,
                });
            }
        });
    },

    // zoneCreate creates a new zone on the server
    zoneCreate: function(zoneJson, callback) {
        $.ajax({
            url: BASE + '/api/v1/zones',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(zoneJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    status: status,
                    validationErrors: errors
                });
            }
        });
    },

    // zoneUpdate updates a zone on the server with the new values
    zoneUpdate: function(zoneJson, callback) {
        $.ajax({
            url: BASE + '/api/v1/zones/' + zoneJson.id,
            type: 'PUT',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(zoneJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                var errors = (xhr.responseJSON || {}).errors;
                callback({
                    err: err,
                    xhr: xhr,
                    validationErrors: errors
                });
            }
        });        
    },

    // zoneSetLevel sets the level of a zone.
    // cmd -> 'turnOn | turnOff | setLevel
    zoneSetLevel: function(zoneId, cmd, value, r, g, b, callback) {
        $.ajax({
            url: BASE + '/api/v1/zones/' + zoneId + '/level',
            type: 'PUT',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify({
                cmd: cmd,
                value: value,
                r: r,
                g: g,
                b: b
            }),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback(err);
            }
        });
    },

    // buttonLoadAll loads all of the buttons in the system
    buttonLoadAll: function(callback) {
        $.ajax({
            url: BASE + '/api/v1/buttons',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            }.bind(this),
            error: function(xhr, status, err) {
                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }.bind(this)
        });
    },

    // discoverersList lists all of the discoverers
    discoverersList: function(callback) {
        $.ajax({
            url: BASE + '/api/v1/discovery/discoverers',
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err
                });
            }
        });
    },

    // discovererScanDevices scans the local network for specific devices
    discovererScanDevices: function(discovererId, callback) {
        $.ajax({
            url: BASE + '/api/v1/discovery/discoverers/' + discovererId,
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err
                });
            }
        });
    },

    // discovererFromString reads the config string and returns goHOME relevant devices
    discovererFromString: function(discovererId, configStr, callback) {
        $.ajax({
            url: BASE + '/api/v1/discovery/discoverers/' + discovererId,
            dataType: 'json',
            type: 'POST',
            data: JSON.stringify(configStr),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback({
                    err: err
                });
            }
        });
    }

};
module.exports = API;
