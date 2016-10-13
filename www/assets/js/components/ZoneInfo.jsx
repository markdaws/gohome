var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var SaveBtn = require('./SaveBtn.jsx');
var DevicePicker = require('./DevicePicker.jsx');
var ZoneOutputPicker = require('./ZoneOutputPicker.jsx');
var ZoneTypePicker = require('./ZoneTypePicker.jsx');

var ZoneInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            clientId: this.props.clientId,
            name: this.props.name,
            description: this.props.description,
            address: this.props.address,
            deviceId: this.props.deviceId,
            type: this.props.type,
            output: this.props.output,
            errors: null,
        }
    },

    toJson: function() {
        var s = this.state
        return {
            clientId: s.clientId,
            name: s.name,
            description: s.description,
            address: s.address,
            deviceId: s.deviceId,
            type: s.type,
            output: s.output,
        }
    },

    setErrors: function(errors) {
        this.setState({ errors: errors });
    },

    _changed: function(evt) {
        this.props.changed && this.props.changed();
        this.changed(evt);
    },
    
    devicePickerChanged: function(deviceId) {
        this.setState({ deviceId: deviceId });
    },

    typeChanged: function(type) {
        this.setState({ type: type });
    },

    outputChanged: function(output) {
        this.setState({ output: output });
    },

    render: function() {
        return (
            <div className="cmp-ZoneInfo well">
              <div className={this.addErr('form-group', 'name')}>
                <label className="control-label" htmlFor={this.uid('name')}>Name*</label>
                <input
                    value={this.state.name}
                    data-statepath="name"
                    onChange={this._changed}
                    className="name form-control"
                    type="text"
                    id={this.uid('name')} />
                {this.errMsg('name')}
              </div>
              <div className={this.addErr("form-group", 'description')}>
                <label className="control-label" htmlFor={this.uid("description")}>Description</label>
                <input
                    value={this.state.description}
                    data-statepath="description"
                    onChange={this._changed}
                    className="description form-control"
                    type="text"
                    id={this.uid("description")} />
                {this.errMsg('description')}
              </div>
              <div className={this.addErr("form-group", "address")}>
                <label className="control-label" htmlFor={this.uid("address")}>Address</label>
                <input
                    value={this.state.address}
                    data-statepath="address"
                    onChange={this._changed}
                    className="address form-control"
                    type="text"
                    id={this.uid("address")} />
                {this.errMsg('address')}
              </div>
              <div className={this.addErr("form-group", "deviceId")}>
                <label className="control-label" htmlFor={this.uid("deviceId")}>Device*</label>
                <DevicePicker
                    disabled={this.isReadOnly("deviceId")}
                    defaultId={this.props.deviceId}
                    devices={this.props.devices}
                    changed={this.devicePickerChanged}/>
                {this.errMsg('deviceId')}
              </div>
              <div className={this.addErr("form-group", "type")}>
                <label className="control-label" htmlFor={this.uid("type")}>Type*</label>
                <ZoneTypePicker type={this.props.type} changed={this.typeChanged}/>
                {this.errMsg('type')}
              </div>
              <div className={this.addErr("form-group", "output")}>
                <label className="control-label" htmlFor={this.uid("output")}>Output*</label>
                <ZoneOutputPicker output={this.props.output} changed={this.outputChanged}/>
                {this.errMsg('output')}
              </div>
            </div>
        );
    }
});
module.exports = ZoneInfo;
