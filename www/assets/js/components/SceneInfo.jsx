var React = require('react');
var SaveBtn = require('./SaveBtn.jsx');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var CommandInfo = require('./CommandInfo.jsx');
var CommandTypePicker = require('./CommandTypePicker.jsx');

var SceneInfo = React.createClass({
    mixins: [InputValidationMixin, UniqueIdMixin],

    getDefaultProps: function() {
        return {
            buttons: []
        };
    },
    
    getInitialState: function() {
        return {
            id: this.props.scene.id || '',
            name: this.props.scene.name || '',
            address: this.props.scene.address || '',
            managed: (this.props.scene.managed == undefined) ? true : this.props.scene.managed,
            commands: this.props.scene.commands || [],
            zones: this.props.zones || [],
            scenes: this.props.scenes || [],
            dirty: false
        };
    },

    componentWillReceiveProps: function(nextProps) {
        //Needed?
        if (nextProps.zones) {
            this.setState({ zones: nextProps.zones });
        }
        if (nextProps.scenes) {
            this.setState({ scenes: nextProps.scenes });
        }
    },
    
    toJson: function() {
        var s = this.state;
        return {
            id: this.state.id,
            name: this.state.name,
            address: this.state.address,
            managed: this.state.managed,
        };
    },

    deleteScene: function() {
        var self = this;
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + this.state.id,
            type: 'DELETE',
            cache: false,
            success: function(data) {
                self.props.onDestroy(self.state.id);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err);
            }.bind(this)
        });
    },

    saveCommand: function(cmd, callback) {
        var self = this;
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + this.state.id + '/commands',
            type: 'POST',
            dataType: 'json',
            data: JSON.stringify(cmd),
            cache: false,
            success: function(data) {
                console.log('saved command');
                callback();
            },
            error: function(xhr, status, err) {
                var errors = (JSON.parse(xhr.responseText) || {}).errors;
                callback(errors);
            }
        });
    },
    
    deleteCommand: function(cmdIndex, isNewCmd, callback) {
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
        var cmds = this.state.commands;
        cmds.push({ isNew: true, type: cmdType, attributes: {} });
        this.setState({ commands: cmds });
    },

    save: function() {
        var saveBtn = this.refs.saveBtn;
        saveBtn.saving();

        this.setState({ errors: null });
        var self = this;
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + this.state.id,
            type: 'PUT',
            dataType: 'json',
            data: JSON.stringify(this.toJson()),
            cache: false,
            success: function(data) {
                self.setState({ dirty: false });
            },
            error: function(xhr, status, err) {
                var errors = (JSON.parse(xhr.responseText) || {}).errors;
                self.setState({ errors: errors });
                saveBtn.failure();
            }
        });        
    },
    
    render: function() {
        var commands
        //TODO: remove
        this.state.managed = true;
        var self = this;
        if (this.state.managed) {
            var cmdIndex = 0;

            commands = this.state.commands.map(function(command) {
                // This isn't a great idea for react, but we don't have anything
                // that can be used as a key since commands don't have ids, will take
                // the perf hit for now
                var key = Math.random();
                var info = (
                    <CommandInfo
                      isNew={command.isNew}
                      key={key}
                      index={cmdIndex}
                      onSave={self.saveCommand}
                      onDelete={self.deleteCommand}
                      scenes={self.props.scenes}
                      zones={self.props.zones}
                      buttons={self.props.buttons}
                      command={command} />
                    );
                cmdIndex++;
                return info;
            });
        } else {
            commands = <p>This is an unmanaged scene. The scene is controlled by a 3rd party device so we can&apos;t show the individual commands it will execute. To modify the scene you will need to use the app provided with the 3rd party device.</p>
        }

        var saveBtn;
        if (this.state.dirty) {
            saveBtn = (
                <div className="pull-right">
                  <SaveBtn text="Save" ref="saveBtn" clicked={this.save} />
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
                  onChange={this.changed}
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
                  onChange={this.changed}
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
                  onChange={this.changed}
                  className="address form-control"
                  type="text"
                  id={this.uid("address")}/>
                {this.errMsg("address")}
              </div>
              <div className="clearfix">
                <a data-toggle="collapse" href={"#" + this.uid("commands")}>
                  Toggle Info
                  <i className="glyphicon glyphicon-menu-down"></i>
                </a>
                {saveBtn}
              </div>
              <div className="collapse commands" id={this.uid("commands")}>
                <h3>Commands</h3>
                {commands}
                Add Command: <CommandTypePicker changed={this.commandTypeChanged}/>
              </div>
            </div>
        );
    }
});
module.exports = SceneInfo;