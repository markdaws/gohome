var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var DevicePicker = require('./DevicePicker.jsx');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SensorInfo',
    prefix: 'b-'
});
require('../../css/components/SensorInfo.less')

var SensorInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            clientId: this.props.clientId,
            name: this.props.name,
            description: this.props.description,
            address: this.props.address,
            deviceId: this.props.deviceId,
            errors: this.props.errors,
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
            attr: this.props.attr,
        }
    },

    setErrors: function(errors) {
        this.setState({ errors: errors });
    },

    _changed: function(evt) {
        this.changed(evt, function(){
            this.props.changed && this.props.changed(this);
        }.bind(this));
    },
    
    devicePickerChanged: function(deviceId) {
        this.setState({ deviceId: deviceId }, function() {
            this.props.changed && this.props.changed(this);
        }.bind(this));
    },

    render: function() {
        return (
            <div {...classes('', '', 'well well-sm')}>
              <div className={this.addErr('form-group', 'name')}>
                <label {...classes('label', '', 'control-label')} htmlFor={this.uid('name')}>Name*</label>
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
                <label {...classes('label', '', 'control-label')} htmlFor={this.uid("description")}>Description</label>
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
                <label {...classes('label', '', 'control-label')} htmlFor={this.uid("address")}>Address</label>
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
                <label {...classes('label', '', 'control-label')} htmlFor={this.uid("deviceId")}>Device*</label>
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
