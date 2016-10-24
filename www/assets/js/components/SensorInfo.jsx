var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var DevicePicker = require('./DevicePicker.jsx');

var SensorInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            clientId: this.props.clientId,
            name: this.props.name,
            description: this.props.description,
            address: this.props.address,
            deviceId: this.props.deviceId,
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

    render: function() {
        return (
            <div className="cmp-DeviceInfo well">
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
            </div>
        );
    }
});
module.exports = SensorInfo;
