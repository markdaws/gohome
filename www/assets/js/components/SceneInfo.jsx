var React = require('react');
var ReactRedux = require('react-redux');
var SaveBtn = require('./SaveBtn.jsx');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var CommandInfo = require('./CommandInfo.jsx');
var CommandTypePicker = require('./CommandTypePicker.jsx');
var SceneActions = require('../actions/SceneActions.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SceneInfo',
    prefix: 'b-'
});
require('../../css/components/SceneInfo.less')

var SceneInfo = React.createClass({
    mixins: [InputValidationMixin, UniqueIdMixin],

    getDefaultProps: function() {
        return {
            //TODO: remove
            buttons: []
        };
    },

    getInitialState: function() {
        return {
            id: this.props.scene.id || '',
            name: this.props.scene.name || '',
            address: this.props.scene.address || '',
            managed: (this.props.scene.managed == undefined) ? true : this.props.scene.managed,
            errors: this.props.errors,

            // true if the object has been modified
            dirty: false,

            saveButtonStatus: this.saveStatus
        };
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.errors) {
            this.setState({ errors: nextProps.errors });
        }
        if (nextProps.saveStatus) {
            this.setState({ saveButtonStatus: nextProps.saveStatus });
        }
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
        var self = this;

        //TODO: Broken
        alert('this is broken, need to know if scene is new vs on server');
        if (this.state.id === '') {
            //TODO: Update state on save, not dirty, has id not clientId
            this.props.saveScene(this.toJson());
        } else {
            //TODO: Verify state correct after successfully updated
            this.props.updateScene(this.toJson());
        }
    },

    deleteScene: function() {
        this.props.deleteScene(this.state.id);
    },

    commandTypeChanged: function(cmdType) {
        this.props.addCommand(this.state.id, cmdType);
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
                commandNodes = <p>To add commands, first save the scene.</p>
            } else {
                var commands = this.props.scene.commands || [];
                commandNodes = commands.map(function(command) {
                    //TODO: We need to give commands an ID on the server so we can have a proper index
                    var key = Math.random();
                    var info = (
                        <CommandInfo
                            isNew={command.isNew}
                            scene={self.props.scene}
                            key={key}
                            index={cmdIndex}
                            scenes={self.props.scenes}
                            zones={self.props.zones}
                            buttons={self.props.buttons}
                            command={command} />
                    );
                    cmdIndex++;
                    return info;
                });
                commandNodes = (
                    <div>
                        {commandNodes}
                        Add Command: <CommandTypePicker changed={this.commandTypeChanged}/>
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
                        Edit Commands
                        <i {...classes('down-arrow', '', 'glyphicon glyphicon-menu-down')}></i>
                    </a>
                    {saveBtn}
                </div>
                <div {...classes('commands', '', 'collapse')} id={this.uid("commands")}>
                    <h3 {...classes('command-header')}>Commands</h3>
                    {commandNodes}
                </div>
            </div>
        );
    }
});

module.exports = SceneInfo;
