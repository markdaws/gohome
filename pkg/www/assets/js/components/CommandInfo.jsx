var React = require('react');
var ReactRedux = require('react-redux');
var SceneSetCommand = require('./SceneSetCommand.jsx');
var SaveBtn = require('./SaveBtn.jsx');
var Feature = require('../feature.js');
var FeatureSetAttrsCommand = require('./FeatureSetAttrsCommand.jsx');
var Api = require('../utils/API.js');
var Constants = require('../constants.js');
var SceneActions = require('../actions/SceneActions.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'CommandInfo',
    prefix: 'b-'
});
require('../../css/components/CommandInfo.less')

var CommandInfo = React.createClass({
    getInitialState: function() {
        return {
            saveButtonStatus: '',
            commandModified: false
        };
    },

    deleteCommand: function() {
        this.props.deleteCommand(
            this.props.scene.id,
            this.props.command.id,
            this.props.command.clientId
        );
    },

    saveCommand: function() {
        this.setState({ errors: [] });

        var data;
        switch (this.props.command.type) {
            case 'featureSetAttrs':
                var featureCmp = this.refs.featureCmp;
                var settings = featureCmp.getSettings();
                data = {
                    type: 'featureSetAttrs',
                    attributes: {
                        id: settings.feature.id,
                        type: settings.feature.type,
                        attrs: settings.modifiedAttrs
                    }
                };
                break;

            case 'sceneSet':
                var sceneCmp = this.refs.sceneCmp;
                data = sceneCmp.toJson();
                break;

            default:
                console.error('unknown command type: ' + this.props.command.type);
        }

        Api.sceneSaveCommand(
            this.props.scene.id,
            data,
            function(err, data) {
                if (err) {
                    this.setState({
                        saveButtonStatus: 'error'
                    });
                    return;
                }

                this.setState({ saveButtonStatus: 'success' });
                this.props.savedCommand(this.props.scene.id, this.props.command.clientId, data);
            }.bind(this)
        );
    },

    commandChanged: function() {
        this.setState({ commandModified: true });
    },

    render: function() {
        var command = this.props.command;
        var saveBtn
        if (this.props.command.clientId && this.state.commandModified) {
            saveBtn = (
                <SaveBtn
                    text="Save"
                    status={this.state.saveButtonStatus}
                    clicked={this.saveCommand} />
            );
        }

        var uiCmd;
        switch (command.type) {
            case 'featureSetAttrs':
                uiCmd = (
                    <FeatureSetAttrsCommand
                        ref='featureCmp'
                        command={command}
                        onAttrChanged={this.commandChanged}
                        devices={this.props.devices}/>
                );
                break;

            case 'sceneSet':
                uiCmd = (
                    <SceneSetCommand
                        ref="sceneCmp"
                        disabled={command.id}
                        parentSceneId={this.props.scene.id}
                        scenes={this.props.scenes}
                        onChanged={this.commandChanged}
                        command={command}/>
                );
                break;

            default:
                console.error('unknown command type: ' + command.type);
                break;
        }

        return (
            <div {...classes('','', 'well well-sm clearfix')}>
                <button {...classes('btn-delete', '', 'btn btn-link pull-right')} onClick={this.deleteCommand}>
                    <i>&#x2717;</i>
                </button>
                {uiCmd}
                <div {...classes('save-btn', '', 'pull-right')}>
                    {saveBtn}
                    <div style={{clear:"both"}}></div>
                </div>
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        savedCommand: function(sceneId, cmdClientId, cmdData) {
            dispatch(SceneActions.savedCommand(sceneId, cmdClientId, cmdData));
        },
        deleteCommand: function(sceneId, cmdId, cmdClientId) {
            dispatch(SceneActions.deleteCommand(sceneId, cmdId, cmdClientId));
        }
    }
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(CommandInfo);
