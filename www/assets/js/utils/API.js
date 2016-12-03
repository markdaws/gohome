var BASE = '//' + window.location.hostname + ':5000';

/*
 API provides helper methods to access all of the gohome REST APIs
*/
var API = {
    //TODO: Use a client side router and middleware

    // setSID sets the session ID, this is needed to call an API
    setSID: function(sid) {
        this.SID = sid;
    },

    checkErr: function(xhr) {
        if (xhr.status === 401) {
            // unauthorized, in this case the user has an invalid sid cookie, delete
            // it and reload the app which will take the user back to the login page
            document.cookie = 'sid=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            window.location = '/';
            return true;
        }

        return false;
    },

    // url builds a valid url to call an API
    url: function(url) {
        return BASE + url + '?sid=' + this.SID;
    },

    monitorGroups: function(monitorId) {
        var conn = new WebSocket('ws:' + this.url('/api/v1/monitor/groups/' + monitorId));
        return conn;
    },

    deviceLoadAll: function(callback) {
        $.ajax({
            url: this.url('/api/v1/devices'),
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }.bind(this)
        });
    },

    deviceCreate: function(deviceJson, callback) {
        $.ajax({
            url: this.url('/api/v1/devices'),
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(deviceJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    deviceUpdate: function(deviceJson, callback) {
        $.ajax({
            url: this.url('/api/v1/devices/' + deviceJson.id),
            type: 'PUT',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(deviceJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    // Deletes a device on the server
    deviceDestroy: function(id, callback) {
        $.ajax({
            url: this.url('/api/v1/devices/' + id),
            type: 'DELETE',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({
                    err: err,
                    status: status,
                    xhr: xhr
                });
            }.bind(this)
        });
    },

    deviceSetFeatureAttrs: function(deviceId, featureId, attrs, callback) {
        $.ajax({
            url: this.url('/api/v1/devices/' + deviceId + '/features/' + featureId + '/apply'),
            type: 'PUT',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(attrs),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    // sceneActivate actives the specified scene
    sceneActivate: function(sceneId, callback) {
        $.ajax({
            url: this.url('/api/v1/scenes/active'),
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify({ id: sceneId }),
            success: function(data) {
                callback(null, data);
            }.bind(this),
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({ err: err });
            }.bind(this)
        });
    },

    // sceneLoadAll loads all of the scenes from the backing store
    sceneLoadAll: function(callback) {
        $.ajax({
            url: this.url('/api/v1/scenes'),
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({
                    err: err,
                    status: status,
                    xhr: xhr,
                });
            }.bind(this)
        });
    },

    // sceneCreate creates a new scene in the backing store
    sceneCreate: function(scene, callback) {
        $.ajax({
            url: this.url('/api/v1/scenes'),
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(scene),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    // sceneUpdate updates fields of an existing scene
    sceneUpdate: function(scene, callback) {
        $.ajax({
            url: this.url('/api/v1/scenes/' + scene.id),
            type: 'PUT',
            dataType: 'json',
            data: JSON.stringify(scene),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    // sceneDestroy deletes the scene with the specified ID from the backing store
    sceneDestroy: function(id, callback) {
        $.ajax({
            url: this.url('/api/v1/scenes/' + id),
            type: 'DELETE',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({
                    err: err,
                    status: status,
                    xhr: xhr
                });
            }.bind(this)
        });
    },

    sceneSaveCommand: function(sceneId, cmd, callback) {
        $.ajax({
            url: this.url('/api/v1/scenes/' + sceneId + '/commands'),
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(cmd),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    sceneDeleteCommand: function(sceneId, cmdIndex, callback) {
        $.ajax({
            url: this.url('/api/v1/scenes/' + sceneId + '/commands/' + cmdIndex),
            type: 'DELETE',
            dataType: 'json',
            data: {},
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }.bind(this)
        });
    },

    // monitorSubscribe requests to subscribe to sensor and zone changes
    monitorSubscribe: function(group, callback) {
        $.ajax({
            url: this.url('/api/v1/monitor/groups'),
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(group),
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({
                    err: err,
                    xhr: xhr,
                    status: status
                });
            }.bind(this)
        });
    },

    // monitorSubscribeRenew renews the subscription to the monitor group, this extends
    // the timeout for the group before the server stops sending updates
    monitorSubscribeRenew: function(monitorId, callback) {
        $.ajax({
            url: this.url('/api/v1/monitor/groups/' + monitorId),
            type: 'PUT',
            dataType: 'json',
            data: null,
            cache: false,
            success: function(data) {
                if (callback) callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                if (callback) {
                    callback({
                        err: err,
                        xhr: xhr,
                        status: status
                    });
                }
            }.bind(this)
        });
    },

    // monitorUnsubscribe unsubscribe the specified monitor id so the client will
    // no longer receive updates when values associated with it change
    monitorUnsubscribe: function(monitorId, callback) {
        $.ajax({
            url: this.url('/api/v1/monitor/groups/' + monitorId),
            type: 'DELETE',
            dataType: 'json',
            data: null,
            cache: false,
            success: function(data) {
                if (callback) callback(null, data);
            },
            error: function(xhr, status, err) {
                if (callback) {
                    callback({
                        err: err,
                        xhr: xhr,
                        status: status
                    });
                }
            }.bind(this)
        });
    },

    // featureUpdate updates a feature on the server with the new values
    featureUpdate: function(deviceId, featureId, featureJson, callback) {
        $.ajax({
            url: this.url('/api/v1/devices/' + deviceId + '/features/' + featureId),
            type: 'PUT',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(featureJson),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    // discoverersList lists all of the discoverers
    discoverersList: function(callback) {
        $.ajax({
            url: this.url('/api/v1/discovery/discoverers'),
            dataType: 'json',
            cache: false,
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({
                    err: err
                });
            }.bind(this)
        });
    },

    // discovererScanDevices scans the local network for specific devices
    discovererScanDevices: function(discovererId, uiFields, callback) {
        $.ajax({
            url: this.url('/api/v1/discovery/discoverers/' + discovererId),
            dataType: 'json',
            cache: false,
            type: 'POST',
            data: JSON.stringify(uiFields),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback(xhr.responseJSON.err);
            }.bind(this)
        });
    },

    sessionCreate: function(login, password, callback) {
        // NOTE: This api lives on the WWW server, so we get a session cookie set on the
        // WWW domain, vs. this being on the API domain
        $.ajax({
            url: '//' + window.location.host + '/api/v1/users/' + login + '/sessions',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify({ password: password }),
            success: function(data) {
                callback(null, data);
            },
            error: function(xhr, status, err) {
                if (this.checkErr(xhr)) {
                    return;
                }

                callback({});
            }.bind(this)
        });
    }
};
module.exports = API;
