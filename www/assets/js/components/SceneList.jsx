var React = require('react');
var ReactDOM = require('react-dom');
var Redux = require('redux');
var ReactRedux = require('react-redux');
var Scene = require('./Scene.jsx');
var SceneInfo = require('./SceneInfo.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
//var SceneStore = require('../stores/SceneStore.js');
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
        //SceneStore.addChangeListener(this._onChange);
        //SceneActions.loadAll();
        this.props.loadAllScenes();

        //TODO: Enable as part of a mode
        //var el = ReactDOM.findDOMNode(this).getElementsByClassName('sceneList')[0];
        //Sortable.create(el);
    },

    componentWillUnmount: function() {
        //SceneStore.removeChangeListener(this._onChange);
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
        var scenes = this.props.scenes.items;//SceneStore.getAll();
        if (this.state.editMode) {
            //TODO: ??
            var newScene = null;//SceneStore.getNewScene();

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
                  <button className="btn btn-primary btnNew pull-left" onClick={this.props.newScene}>New Scene</button>
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

function mapStateToProps(state) {
    return {
        scenes: state.scenes
    }
}

function mapDispatchToProps(dispatch) {
    return {
        newScene: function() {
            dispatch(SceneActions.newClient());
        },
        loadAllScenes: function() {
            dispatch(SceneActions.loadAll());
        }
    }
}

var SceneListContainer = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(SceneList);
module.exports = SceneListContainer;
