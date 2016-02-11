var React = require('react');
var ReactDOM = require('react-dom');
var Scene = require('./Scene.jsx');
var SceneInfo = require('./SceneInfo.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var SceneStore = require('../stores/SceneStore.js');
var SceneActions = require('../actions/SceneActions.js');

var SceneList = React.createClass({
    mixins: [UniqueIdMixin],
    
    getInitialState: function() {
        return {
            editMode: false,
            //TODO: remove
            zones: []
        };
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.zones) {
            this.setState({ zones: nextProps.zones });
        }
    },
    
    componentDidMount: function() {
        SceneStore.addChangeListener(this._onChange);
        SceneActions.loadAll();

        //TODO: Enable as part of a mode
        //var el = ReactDOM.findDOMNode(this).getElementsByClassName('sceneList')[0];
        //Sortable.create(el);
    },

    componentWillUnmount: function() {
        SceneStore.removeChangeListener(this._onChange);
    },

    _onChange: function() {
        this.forceUpdate();
    },

    edit: function() {
        this.setState({ editMode: true });
    },

    endEdit: function() {
        this.setState({ editMode: false });
    },

    newScene: function() {
        SceneActions.newClient();
    },

    render: function() {
        var body;
        var btns;
        var scenes = SceneStore.getAll();
        if (this.state.editMode) {
            var newScene = SceneStore.getNewScene();

            // If the user is in the process of creating a new scene we append the
            // current new scene object to the front of the list
            if (newScene) {
                scenes = scenes.unshift(newScene);
            }

            var self = this;
            body = scenes.map(function(scene) {
                return (
                    <SceneInfo
                      zones={self.state.zones}
                      buttons={self.props.buttons}
                      scene={scene}
                      readOnlyFields="id"
                      key={scene.id || scene.clientId} />
                );
            });
            btns = (
                <div className="clearfix buttonWrapper">
                  <button className="btn btn-primary btnNew pull-left" onClick={this.newScene}>New Scene</button>
                  <button className="btn btn-success btnDone pull-right" onClick={this.endEdit}>Done</button>
                </div>
            );
        } else {

            console.log(scenes)
            body = scenes.map(function(scene) {
                return (
                    <Scene scene={scene} key={scene.id}/>
                );
            });
            btns = (
                <div className="clearfix buttonWrapper">
                  <button className="btn btn-primary btnEdit pull-right" onClick={this.edit}>Edit</button>
                </div>
            );
        }
        
        return (
            <div className="cmp-SceneList">
              {btns}
              {body}
            </div>
        );
    }
});
module.exports = SceneList;