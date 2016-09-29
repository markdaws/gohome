var React = require('react');
var SaveBtn = require('./SaveBtn.jsx');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var CommandInfo = require('./CommandInfo.jsx');
var CommandTypePicker = require('./CommandTypePicker.jsx');
var SceneActions = require('../actions/SceneActions.js');

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
            clientId: this.props.scene.clientId,
            name: this.props.scene.name || '',
            address: this.props.scene.address || '',
            managed: (this.props.scene.managed == undefined) ? true : this.props.scene.managed,
            errors: this.props.errors,
            //TODO: Needed?, turn to props
            commands: this.props.scene.commands || [],

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
        var s = this.state;
        return {
            id: this.state.id,
            name: this.state.name,
            address: this.state.address,
            managed: this.state.managed,
            clientId: this.state.clientId,
        };
    },

    saveScene: function() {
        this.setState({ errors: null });
        var self = this;

        if (this.state.id === '') {
            //TODO: Update state on save, not dirty, has id not clientId
            this.props.saveScene(this.toJson());
        } else {
            //TODO: Verify state correct after successfully updated
            this.props.updateScene(this.toJson());
        }
    },

    deleteScene: function() {
        this.props.deleteScene(this.state.clientId, this.state.id);
    },

    addCommand: function(cmd, callback) {
        this.props.addCommand(this.state.id, cmd);
        //TODO: remove

        /*
           var self = this;
           $.ajax({
           url: '/api/v1/systems/123/scenes/' + this.state.id + '/commands',
           type: 'POST',
           dataType: 'json',
           data: JSON.stringify(cmd),
           cache: false,
           success: function(data) {
           callback();
           },
           error: function(xhr, status, err) {
           var errors = (JSON.parse(xhr.responseText) || {}).errors;
           callback(errors);
           }
           });*/
    },

    deleteCommand: function(cmdIndex, isNewCmd, callback) {
        //TODO: redux
        var self = this;
        
        if (isNewCmd) {
            var commands = self.state.commands.filter(function(cmd, index) {
                return index != cmdIndex;
            });
            self.setState({ commands: commands });
            return;
        }

        $.ajax({
            url: '/api/v1/systems/123/scenes/' + this.state.id + '/commands/' + cmdIndex,
            type: 'DELETE',
            cache: false,
            success: function(data) {
                var commands = self.state.commands.filter(function(cmd, index) {
                    return index != cmdIndex;
                });
                self.setState({ commands: commands });
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err);
                callback(err);
            }.bind(this)
        });
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
                            key={key}
                            index={cmdIndex}
                            onSave={self.addCommand}
                            onDelete={self.deleteCommand}
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
            commandNodes = <p>This is an unmanaged scene. The scene is controlled by a 3rd party device so we can&apos;t show the individual commands it will execute. To modify the scene you will need to use the app provided with the 3rd party device.</p>
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
            <div className="cmp-SceneInfo well well-sm">
                <button className="btn btn-link btnDelete pull-right" onClick={this.deleteScene}>
                    <i className="glyphicon glyphicon-trash"></i>
                </button>
                <div className={this.addErr("form-group", "name")}>
                    <label className="control-label" htmlFor={this.uid("name")}>Name*</label>
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
                    <label className="control-label" htmlFor={this.uid("id")}>ID</label>
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
                    <label className="control-label" htmlFor={this.uid("address")}>Address</label>
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
                        <i className="glyphicon glyphicon-menu-down"></i>
                    </a>
                    {saveBtn}
                </div>
                <div className="collapse commands" id={this.uid("commands")}>
                    <h3>Commands</h3>
                    {commandNodes}
                </div>
            </div>
        );
    }
});
module.exports = SceneInfo;
