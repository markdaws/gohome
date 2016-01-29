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
        //TODO:
    },

    save: function() {
        var saveBtn = this.refs.saveBtn;
        saveBtn.saving();
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
            uiCmd = <ZoneSetLevelCommand zones={this.props.zones} command={command} />;
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