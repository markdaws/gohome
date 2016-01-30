var React = require('react');
var ZoneSetLevelCommand = require('./ZoneSetLevelCommand.jsx');
var SaveBtn = require('./SaveBtn.jsx');

var CommandInfo = React.createClass({
    getInitialState: function() {
        return {
            command: this.props.command
        }
    },

    deleteCommand: function() {
        this.props.onDelete(this.props.index, function(err) {
            console.log('I was deleted: ' + err);
            // TODO: If there is an error then the delete button should
            // show an error state ... 
        });
    },

    save: function() {
        var saveBtn = this.refs.saveBtn;
        saveBtn.saving();

        var cmd = this.refs.cmd.toJson();
        this.props.onSave(cmd, function(err) {
            console.log('save cb: ' + err);
        });
    },
    
    render: function() {
        var self = this;
        var command = this.state.command;
        var uiCmd;
        switch (command.type) {
        case 'buttonPress':
            //TODO:
            break;
        case 'buttonRelease':
            //TODO:
            break;
        case 'zoneSetLevel':
            uiCmd = <ZoneSetLevelCommand ref="cmd" zones={this.props.zones} command={command} />;
            break;
        case 'sceneSet':
            break;
        default:
            console.error('unknown command type: ' + command.type);
        }
        return (
            <div className="cmp-CommandInfo well clearfix">
              {uiCmd}
              <button className="btn btn-danger btnDelete pull-right" onClick={this.deleteCommand}>Delete</button>
              <SaveBtn text="Save" ref="saveBtn" clicked={this.save}/>
            </div>
        );
    }
});
module.exports = CommandInfo;