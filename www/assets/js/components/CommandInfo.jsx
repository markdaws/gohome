var React = require('react');
var ZoneSetLevelCommand = require('./ZoneSetLevelCommand.jsx');
var SceneSetCommand = require('./SceneSetCommand.jsx');
var SaveBtn = require('./SaveBtn.jsx');
var ButtonPressCommand = require('./ButtonPressCommand.jsx');
var ButtonReleaseCommand = require('./ButtonReleaseCommand.jsx');

var CommandInfo = React.createClass({
    getInitialState: function() {
        return {
            command: this.props.command,
            isNew: this.props.isNew
        }
    },

    deleteCommand: function() {
        this.props.onDelete(this.props.index, this.state.isNew, function(err) {
            console.log('I was deleted: ' + err);
            // TODO: If there is an error then the delete button should
            // show an error state ... 
        });
    },

    save: function() {
        var saveBtn = this.refs.saveBtn;
        saveBtn.saving();

        var cmd = this.refs.cmd;
        var self = this;
        this.props.onSave(cmd.toJson(), function(errors) {
            if (errors) {
                cmd.setErrors(errors);
                saveBtn.failure();
            } else {
                cmd.setErrors(null);
                self.setState({ isNew: false });
            }
        });
    },
    
    render: function() {
        var self = this;
        var command = this.state.command;
        var saveBtn
        if (this.state.isNew) {
            saveBtn = <SaveBtn text="Save" ref="saveBtn" clicked={this.save}/>
        }
        
        var uiCmd;
        switch (command.type) {
        case 'buttonPress':
            uiCmd = <ButtonPressCommand ref="cmd" buttons={this.props.buttons} command={command}/>;
            break;
        case 'buttonRelease':
            uiCmd = <ButtonReleaseCommand ref="cmd" buttons={this.props.buttons} command={command}/>;
            break;
        case 'zoneSetLevel':
            uiCmd = <ZoneSetLevelCommand ref="cmd" zones={this.props.zones} command={command} />;
            break;
        case 'sceneSet':
            uiCmd = <SceneSetCommand ref="cmd" scenes={this.props.scenes} command={command} />;
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
module.exports = CommandInfo;