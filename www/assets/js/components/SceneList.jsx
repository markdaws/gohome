var React = require('react');
var ReactRedux = require('react-redux');
var SceneListGridCell = require('./SceneListGridCell.jsx');
var SceneControl = require('./SceneControl.jsx');
var SceneInfo = require('./SceneInfo.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var SceneActions = require('../actions/SceneActions.js');
var Grid = require('./Grid.jsx');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SceneList',
    prefix: 'b-'
});
require('../../css/components/SceneList.less')

var SceneList = React.createClass({
    mixins: [UniqueIdMixin],

    getInitialState: function() {
        return { editMode: false };
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
                saveState = (this.props.scenes.saveState[scene.id] || {});

                return (
                    <div {...classes('scene-info')} key={scene.id}>
                        <SceneInfo
                            zones={this.props.zones}
                            buttons={this.props.buttons}
                            scenes={this.props.scenes.items}
                            scene={scene}
                            readOnlyFields="id"
                            key={scene.id}
                            errors={(saveState.err || {}).validationErrors}
                            saveScene={this.props.saveScene}
                            updateScene={this.props.updateScene}
                            deleteScene={this.props.deleteScene}
                            addCommand={this.props.addCommand}
                            saveStatus={saveState.status} />
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
                    key: scene.id,
                    cell: <SceneListGridCell scene={scene} />,
                    content: <SceneControl scene={scene} key={scene.id}/>
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
        deleteScene: function(id) {
            alert('broken needs client id');
            if (id) {
                dispatch(SceneActions.destroyClient(id));
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
