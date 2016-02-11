var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var Constants = require('../constants/constants.js');
var EventEmitter = require('events').EventEmitter;
var Assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _zones = [];

var ZoneStore = Assign({}, EventEmitter.prototype, {
    addChangeListener: function(cb) {
        this.on(CHANGE_EVENT, cb);
    },

    removeChangeListener: function(cb) {
        this.removeListener(CHANGE_EVENT, cb);
    },

    emitChange: function() {
        this.emit(CHANGE_EVENT);
    },

    init: function(zones) {
        _zones = zones;
        this.emitChange();
    },
    
    getAll: function() {
        return _zones;
    },

    dispatcherIndex: AppDispatcher.register(function(payload) {
        switch(payload.actionType) {
        case Constants.ZONE_LOAD_ALL_RAW:
            ZoneStore.init(payload.raw);
            break;
        }
        return true;
    })
});
module.exports = ZoneStore;
