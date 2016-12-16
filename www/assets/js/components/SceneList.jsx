var React = require('react');
var ReactRedux = require('react-redux');
var SceneListGridCell = require('./SceneListGridCell.jsx');
var SceneControl = require('./SceneControl.jsx');
var SceneInfo = require('./SceneInfo.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var SceneActions = require('../actions/SceneActions.js');
var Grid = require('./Grid.jsx');
var BEMHelper = require('react-bem-helper');
var Feature = require('../feature.js');

var classes = new BEMHelper({
    name: 'SceneList',
    prefix: 'b-'
});
require('../../css/components/SceneList.less')

var SceneList = React.createClass({
    mixins: [UniqueIdMixin],

    getDefaultProps: function() {
        return {
            scenes: [],
            devices: []
        };
    },

    getInitialState: function() {
        return {
            // If we don't have any scenes, then we immediately enter in edit mode
            editMode: this.props.scenes.length === 0
        };
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

        var scenes = this.props.scenes;
        var gridCells = [];
        if (this.state.editMode) {
            body = scenes.map(function(scene) {
                return (
                    <div {...classes('scene-info')} key={scene.id || scene.clientId}>
                        <SceneInfo
                            scenes={this.props.scenes}
                            devices={this.props.devices}
                            scene={scene}
                            readOnlyFields="id"
                            key={scene.id || scene.clientId}
                            createdScene={this.props.createdScene}
                            updatedScene={this.props.updatedScene}
                            deleteScene={this.props.deleteScene}
                            addCommand={this.props.addCommand} />
                    </div>
                );
            }.bind(this));
            btns = (
                <div {...classes('buttons', '', 'clearfix')}>
                    <button className="btn btn-default btnNew pull-left" onClick={this.props.newClientScene}>
                        <i className="fa fa-plus"></i>
                    </button>
                    <button className="btn btn-default btnDone pull-right" onClick={this.endEdit}>
                        <i className="fa fa-times"></i>
                    </button>
                </div>
            );
        } else {

            var gridCells = scenes.map(function(scene) {
                return {
                    key: scene.id || scene.clientId,
                    cell: <SceneListGridCell scene={scene} />,
                    content: <SceneControl scene={scene} key={scene.id || scene.clientId}/>
                };
            });
            btns = (
                <div {...classes('buttons', '', 'clearfix')}>
                    <button className="btn btn-default btnEdit pull-right" onClick={this.edit}>
                        <i className="fa fa-cog" aria-hidden="true"></i>
                    </button>
                </div>
            );

            if (gridCells.length > 0) {
                body = <Grid cells={gridCells} />
            }
        }

        return (
            <div {...classes()}>
                <h2 {...classes('header')}>Scenes</h2>
                {btns}
                {body}
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        newClientScene: function() {
            dispatch(SceneActions.newClient());
        },

        deleteScene: function(id, clientId) {
            if (clientId) {
                dispatch(SceneActions.destroyClient(clientId));
            } else {
                dispatch(SceneActions.destroy(id));
            }
        },

        addCommand: function(sceneId, cmd) {
            dispatch(SceneActions.addCommand(sceneId, cmd));
        },

        createdScene: function(sceneJson, clientId) {
            dispatch(SceneActions.created(sceneJson, clientId));
        },

        updatedScene: function(sceneJson, id) {
            dispatch(SceneActions.updated(sceneJson));
        },
    }
}

module.exports = ReactRedux.connect(null, mapDispatchToProps)(SceneList);
