var React = require('react');
var ReactRedux = require('react-redux');
var ZoneSetLevelCommand = require('./ZoneSetLevelCommand.jsx');
var SceneSetCommand = require('./SceneSetCommand.jsx');
var SaveBtn = require('./SaveBtn.jsx');
var ButtonPressCommand = require('./ButtonPressCommand.jsx');
var ButtonReleaseCommand = require('./ButtonReleaseCommand.jsx');
var Api = require('../utils/API.js');
var Constants = require('../constants.js');
var SceneActions = require('../actions/SceneActions.js');

var CommandInfo = React.createClass({
    deleteCommand: function() {
        //TODO: Show error in UI
        //TODO: Normalize middleware to handle showing error from API, standardize responses
        this.props.deleteCommand(this.props.scene.id, this.props.index, this.props.command.isNew);
    },

    saveCommand: function() {
        var cmd = this.refs.cmd;
        this.setState({ errors: [] });

        var cmdJson = cmd.toJson();
        Api.sceneSaveCommand(this.props.scene.id, cmdJson, function(err, data) {
            if (err) {
                cmd.setErrors(err.validationErrors);
                return;
            }

            this.props.savedCommand(cmdJson, this.props.scene.id, this.props.index);
        }.bind(this));
    },

    render: function() {
        var command = this.props.command;
        var saveBtn
        if (this.props.command.isNew) {
            saveBtn = (
                <SaveBtn
                    text="Save"
                    status=""
                    clicked={this.saveCommand} />
            );
        }

        var uiCmd;
        switch (command.type) {
            case 'buttonPress':
                uiCmd = (<ButtonPressCommand
                             ref="cmd"
                             errors={(command.errors || {}).validationErrors}
                             buttons={this.props.buttons}
                             command={command}/>
                );
                break;
            case 'buttonRelease':
                uiCmd = (<ButtonReleaseCommand
                             ref="cmd"
                             errors={(command.errors || {}).validationErrors}
                             buttons={this.props.buttons}
                             command={command}/>
                );
                break;
            case 'zoneSetLevel':
                uiCmd = (<ZoneSetLevelCommand
                    ref="cmd"
                    errors={(command.errors || {}).validationErrors}
                    zones={this.props.zones}
                    command={command} />
                )
                break;
            case 'sceneSet':
                uiCmd = (<SceneSetCommand
                    ref="cmd"
                    errors={(command.errors || {}).validationErrors}
                    scenes={this.props.scenes}
                    command={command} />
                )
                break;
            default:
                console.error('unknown command type: ' + command.type);
        }
        return (
            <div className="cmp-CommandInfo well well-sm clearfix">
                <button className="btn btn-link btnDelete pull-right" onClick={this.deleteCommand}>
                    <i className="glyphicon glyphicon-trash"></i>
                </button>
                {uiCmd}
                {saveBtn}
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        savedCommand: function(cmdData, sceneId, cmdIndex) {
            dispatch({
                type: Constants.SCENE_COMMAND_SAVE_RAW,
                data: cmdData,
                sceneId: sceneId,
                cmdIndex: cmdIndex });
        },
        deleteCommand: function(sceneId, cmdIndex, isNew) {
            dispatch(SceneActions.deleteCommand(sceneId, cmdIndex, isNew));
        }
    }
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(CommandInfo);
