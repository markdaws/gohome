var React = require('react');
var ReactRedux = require('react-redux');
var SaveBtn = require('./SaveBtn.jsx');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var CommandInfo = require('./CommandInfo.jsx');
var SceneActionPicker = require('./SceneActionPicker.jsx');
var Feature = require('../feature.js');
var SceneActions = require('../actions/SceneActions.js');
var BEMHelper = require('react-bem-helper');
var Api = require('../utils/API.js');
var Uuid = require('uuid');

var classes = new BEMHelper({
    name: 'SceneInfo',
    prefix: 'b-'
});
require('../../css/components/SceneInfo.less')

var SceneInfo = React.createClass({
    mixins: [InputValidationMixin, UniqueIdMixin],

    getInitialState: function() {
        return {
            id: this.props.scene.id || '',
            name: this.props.scene.name || '',
            address: this.props.scene.address || '',
            managed: (this.props.scene.managed == undefined) ? true : this.props.scene.managed,
            errors: this.props.errors,
            saveButtonStatus: '',
            dirty: false,
        };
    },

    toJson: function() {
        return {
            id: this.state.id,
            name: this.state.name,
            address: this.state.address,
            managed: this.state.managed,
        };
    },

    saveScene: function() {
        this.setState({ errors: null });

        if (!this.state.id) {
            Api.sceneCreate(this.toJson(), function(err, data) {
                if (err && !err.validation) {
                    this.setState({saveButtonStatus: 'error' });
                    return;
                } else if (err && err.validation) {
                    this.setState({
                        saveButtonStatus: 'error',
                        errors: err.validation.errors[this.state.id]
                    });
                    return;
                }

                this.setState({ saveButtonStatus: 'success' });
                this.props.createdScene(data, this.props.scene.clientId);
            }.bind(this));
        } else {
            Api.sceneUpdate(this.toJson(), function(err, data) {
                if (err && !err.validation) {
                    this.setState({saveButtonStatus: 'error' });
                    return;
                } else if (err && err.validation) {
                    this.setState({
                        saveButtonStatus: 'error',
                        errors: err.validation.errors[this.state.id]
                    });
                    return;
                }

                this.setState({ saveButtonStatus: 'success' });
                this.props.updatedScene(data, this.props.scene.id);
            }.bind(this));
        }
    },

    deleteScene: function() {
        this.props.deleteScene(this.state.id, this.props.scene.clientId);
    },

    sceneActionPickerChanged: function(actionType) {
        var cmd = {
            clientId: Uuid.v4()
        };

        switch (actionType) {
            case 'sceneSet':
                cmd.type = 'sceneSet';
                cmd.attributes = { };
                break;
            default:
                // All other actions are related to features
                cmd.type = 'featureSetAttrs';
                cmd.attributes = {
                    type: actionType
                };
        }
        this.props.addCommand(this.state.id, cmd);
    },

    _inputChanged: function(evt) {
        this.setState({ saveButtonStatus: ''});

        // Lives in InputValidationMixin
        this.changed(evt);
    },

    render: function() {
        var commandNodes

        //TODO: remove
        this.state.managed = true;
        var self = this;
        if (this.state.managed) {
            var cmdIndex = 0;

            if (this.state.id === '') {
                commandNodes = <p>To add actions, first save the scene.</p>
            } else {
                var commands = this.props.scene.commands || [];
                commandNodes = commands.map(function(command) {
                    var info = (
                        <CommandInfo
                            scene={self.props.scene}
                            key={command.id || command.clientId}
                            index={cmdIndex}
                            devices={self.props.devices}
                            scenes={self.props.scenes}
                            command={command} />
                    );
                    cmdIndex++;
                    return info;
                });

                var excluded = {};
                excluded[Feature.Type.Sensor] = true;

                //TODO: Add back once buttons are supported
                excluded[Feature.Type.Button] = true;
                excluded[Feature.Type.CoolZone] = true;

                commandNodes = (
                    <div>
                        {commandNodes}
                        <div {...classes('feature-picker')}>
                            <SceneActionPicker excluded={excluded} changed={this.sceneActionPickerChanged}/>
                        </div>
                    </div>
                );
            }
        } else {
            commandNodes = <p>The scene is controlled by a 3rd party device so we can&apos;t show the individual commands it will execute. To modify the scene you will need to use the app provided with the 3rd party device.</p>
        }

        var saveBtn;
        if (this.state.dirty) {
            var saveResult;
            saveBtn = (
                <div className="pull-right">
                    <SaveBtn
                        text="Save"
                        status={this.state.saveButtonStatus}
                        clicked={this.saveScene} />
                </div>
            );
        }
        return (
            <div {...classes('', '', 'well well-sm')}>
                <button {...classes('delete', '', 'btn btn-link pull-right')} onClick={this.deleteScene}>
                    <i className="glyphicon glyphicon-trash"></i>
                </button>
                <div className={this.addErr("form-group", "name")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("name")}>Name*</label>
                    <input
                        value={this.state.name}
                        data-statepath="name"
                        onChange={this._inputChanged}
                        className="name form-control"
                        type="text"
                        id={this.uid("name")}/>
                    {this.errMsg("name")}
                </div>
                <div className={this.addErr("form-group", "id")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("id")}>ID</label>
                    <input
                        value={this.state.id}
                        readOnly={this.isReadOnly("id")}
                        data-statepath="id"
                        onChange={this._inputChanged}
                        className="id form-control"
                        type="text"
                        id={this.uid("id")}/>
                    {this.errMsg("id")}
                </div>
                <div className={this.addErr("form-group", "address")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("address")}>Address</label>
                    <input
                        value={this.state.address}
                        data-statepath="address"
                        onChange={this._inputChanged}
                        className="address form-control"
                        type="text"
                        id={this.uid("address")}/>
                    {this.errMsg("address")}
                </div>
                <div className="clearfix">
                    <a data-toggle="collapse" href={"#" + this.uid("commands")}>
                        Edit Actions
                        <i {...classes('down-arrow', '', 'glyphicon glyphicon-menu-down')}></i>
                    </a>
                    {saveBtn}
                </div>
                <div {...classes('commands', '', 'collapse')} id={this.uid("commands")}>
                    <h3 {...classes('command-header')}>Actions</h3>
                    {commandNodes}
                </div>
            </div>
        );
    }
});
module.exports = SceneInfo;
