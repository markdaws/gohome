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
        };
    },

    componentWillReceiveProps: function(nextProps) {
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

        var error;
        //TODO: loading error

        var loading;
        if (this.props.scenes.loading) {
            loading = (
                <div className="spinnerWrapper">
                    <p>Loading Scenes ...</p>
                    <i className="fa fa-spinner fa-spin"></i>
                </div>
            );
        }

        var scenes = this.props.scenes.items;
        if (this.state.editMode) {
            body = scenes.map(function(scene) {
                var saveState;

                // Check for input validation errors from the server
                saveState = (this.props.scenes.saveState[scene.clientId || scene.id] || {});

                return (
                    <SceneInfo
                        zones={this.props.zones}
                        buttons={this.props.buttons}
                        scene={scene}
                        readOnlyFields="id"
                        key={scene.id || scene.clientId}
                        errors={(saveState.err || {}).validationErrors}
                        saveScene={this.props.saveScene}
                        updateScene={this.props.updateScene}
                        deleteScene={this.props.deleteScene}
                        addCommand={this.props.addCommand}
                        saveStatus={saveState.status} />
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
                    <button className="btn btn-default btnEdit pull-right" onClick={this.edit}>
                        <i className="fa fa-cog" aria-hidden="true"></i>
                    </button>
                </div>
            );
        }

        if (loading) {
            btns = null;
            body = null;
        }

        return (
            <div className="cmp-SceneList">
                {error}
                {loading}
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
        saveScene: function(sceneJson) {
            dispatch(SceneActions.create(sceneJson));
        },
        updateScene: function(sceneJson) {
            dispatch(SceneActions.update(sceneJson));
        },
        deleteScene: function(clientId, id) {
            if (clientId) {
                dispatch(SceneActions.destroyClient(clientId));
            } else {
                dispatch(SceneActions.destroy(id));
            }
        },
        addCommand: function(sceneId, cmdType) {
            dispatch(SceneActions.addCommand(sceneId, cmdType));
        },
        saveCommand: function(sceneId, cmd) {
            //TODO:
            dispatch(SceneActions.addCommand(sceneId, cmd));
        }
    }
}

var SceneListContainer = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(SceneList);
module.exports = SceneListContainer;
