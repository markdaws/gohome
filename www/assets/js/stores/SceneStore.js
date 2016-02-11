var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var Constants = require('../constants/constants.js');
var EventEmitter = require('events').EventEmitter;
var Assign = require('object-assign');
var Immutable = require('immutable');

var CHANGE_EVENT = 'change';

var _scenes = Immutable.List();
var _newScene = null;
var _clientId = 0;

var SceneStore = Assign({}, EventEmitter.prototype, {
    addChangeListener: function(cb) {
        this.on(CHANGE_EVENT, cb);
    },

    removeChangeListener: function(cb) {
        this.removeListener(CHANGE_EVENT, cb);
    },

    emitChange: function() {
        this.emit(CHANGE_EVENT);
    },

    init: function(scenes) {
        _scenes = scenes;
        this.emitChange();
    },
    
    getAll: function() {
        return _scenes;
    },

    getNewScene: function() {
        return _newScene;
    },

    dispatcherIndex: AppDispatcher.register(function(payload) {
        console.log(payload);

        switch(payload.actionType) {
        case Constants.SCENE_LOAD_ALL_RAW:
            SceneStore.init(Immutable.List(payload.raw));
            break;
        case Constants.SCENE_NEW_CLIENT:
            _newScene = { clientId: 'scene_cid_' + _clientId + '' };
            ++_clientId;
            SceneStore.emitChange();
            break;
        case Constants.SCENE_CREATE_RAW:
            _scenes = _scenes.push(payload.raw);
            _newScene = null;
            SceneStore.emitChange();
            break;
        case Constants.SCENE_DESTROY_RAW:
            _scenes = _scenes.filter(function(scene) {
                return scene.id !== payload.id;
            });
            SceneStore.emitChange();
            break;
        case Constants.SCENE_UPDATE:
            //TODO:
            break;
        case Constants.SCENE_CREATE_COMMAND:
            //TODO:
            break;
        case Constants.SCENE_DESTROY_COMMAND:
            //TODO:
            break;
        }
        return true;
    })
});
module.exports = SceneStore;
