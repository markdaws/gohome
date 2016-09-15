var React = require('react');
var ReactDOM = require('react-dom');
var Redux = require('redux');
var ReactRedux = require('react-redux');
var Scene = require('./Scene.jsx');
var SceneInfo = require('./SceneInfo.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
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
        this.props.loadAllScenes();

        //TODO: Enable as part of a mode
        //var el = ReactDOM.findDOMNode(this).getElementsByClassName('sceneList')[0];
        //Sortable.create(el);
    },

    componentWillUnmount: function() {
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

    render: function() {
        var body;
        var btns;

        //TODO: What about loading fail?
        //this.props.scenes.loading

        var scenes = this.props.scenes.items;
        if (this.state.editMode) {
            var newScene = this.props.scenes.newSceneInfo;

            // If the user is in the process of creating a new scene we append the
            // current new scene object to the front of the list
            if (newScene) {
                // Since we are modifying the array for rendering, shallow copy array
                scenes = scenes.slice();
                scenes.unshift(newScene.scene);
            }

            body = scenes.map(function(scene) {
                var createResponse;
                return (
                    <SceneInfo
                        zones={this.state.zones}
                        buttons={this.props.buttons}
                        scene={scene}
                        readOnlyFields="id"
                        key={scene.id || scene.clientId}
                        saveScene={this.props.saveScene}
                        saveStatus={(newScene || {}).saveStatus} />
                );
            }.bind(this));
            btns = (
                <div className="clearfix buttonWrapper">
                  <button className="btn btn-primary btnNew pull-left" onClick={this.props.newClientScene}>New Scene</button>
                  <button className="btn btn-success btnDone pull-right" onClick={this.endEdit}>Done</button>
                </div>
            );
        } else {

            body = scenes.map(function(scene) {
                return (
                    <Scene scene={scene} key={scene.id || scene.clientId}/>
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

function mapStateToProps(state) {
    return {
        scenes: state.scenes
    }
}

function mapDispatchToProps(dispatch) {
    return {
        newClientScene: function() {
            dispatch(SceneActions.newClient());
        },
        loadAllScenes: function() {
            dispatch(SceneActions.loadAll());
        },
        saveScene: function(scene) {
            dispatch(SceneActions.create(scene));
        }
    }
}

var SceneListContainer = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(SceneList);
module.exports = SceneListContainer;

//TODO: Hide "New Scene" button after click
