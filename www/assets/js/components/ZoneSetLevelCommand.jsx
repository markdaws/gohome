var React = require('react');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var ZonePicker = require('./ZonePicker.jsx');

var ZoneSetLevelCommand = module.exports = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            level: this.props.command.attributes.Level || 0,
            zoneId: this.props.command.attributes.ZoneID || ''
        }
    },

    toJson: function() {
        return {
            Level: this.state.level,
            ZoneID: this.state.zoneId
        };
    },

    zonePickerChanged: function(zoneId) {
        this.setState({ zoneId: zoneId });
    },
    
    render: function() {
        return (
            <div className="cmp-ZoneSetLevelCommand">
              <h4>Zone Set Level</h4>
              <div className={this.addErr("form-group", "zoneId")}>
                <label className="control-label" htmlFor={this.uid("zoneId")}>Zone</label>
                <ZonePicker changed={this.zonePickerChanged} zones={this.props.zones} zoneId={this.state.zoneId} />
                {this.errMsg("zoneId")}
              </div>
              <div className={this.addErr("form-group", "level")}>
                <label className="control-label" htmlFor={this.uid("level")}>Level</label>
                <input value={this.state.level} data-statepath="level" onChange={this.changed} className="level form-control" type="text" id={this.uid("level")}/>
                {this.errMsg("level")}
              </div>
            </div>
        );
    }
});
module.exports = ZoneSetLevelCommand;