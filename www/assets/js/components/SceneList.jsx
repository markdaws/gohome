var React = require('react');
var ReactDOM = require('react-dom');
var Redux = require('redux');
var ReactRedux = require('react-redux');
var SceneListGridCell = require('./SceneListGridCell.jsx');
var SceneControl = require('./SceneControl.jsx');
var SceneInfo = require('./SceneInfo.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var SceneActions = require('../actions/SceneActions.js');
var Grid = require('./Grid.jsx');

var SceneList = React.createClass({
    mixins: [UniqueIdMixin],

    getInitialState: function() {
        return {
            editMode: false,
        };
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

        var scenes = this.props.scenes.items;
        var gridCells = [];
        if (this.state.editMode) {
            body = scenes.map(function(scene) {
                var saveState;

                // Check for input validation errors from the server
                saveState = (this.props.scenes.saveState[scene.clientId || scene.id] || {});

                return (
                    <SceneInfo
                        zones={this.props.zones}
                        buttons={this.props.buttons}
                        scenes={this.props.scenes.items}
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

            var gridCells = scenes.map(function(scene) {
                return {
                    cell: <SceneListGridCell scene={scene} />,
                    content: <SceneControl scene={scene} key={scene.id || scene.clientId}/>
                };
            });
            btns = (
                <div className="clearfix buttonWrapper">
                    <button className="btn btn-default btnEdit pull-right" onClick={this.edit}>
                        <i className="fa fa-cog" aria-hidden="true"></i>
                    </button>
                </div>
            );

            body = <Grid cells={gridCells} />
        }

        return (
            <div className="cmp-SceneList">
                <h2>Scenes</h2>
                {btns}
                {body}
            </div>
        );
    }
});

function mapStateToProps(state) {
    return { }
}

function mapDispatchToProps(dispatch) {
    return {
        newClientScene: function() {
            dispatch(SceneActions.newClient());
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
        }
    }
}

module.exports = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(SceneList);
