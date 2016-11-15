var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var DevicePicker = require('./DevicePicker.jsx');
var BEMHelper = require('react-bem-helper');
var SaveBtn = require('./SaveBtn.jsx');
var Api = require('../utils/API.js');

var classes = new BEMHelper({
    name: 'SensorInfo',
    prefix: 'b-'
});
require('../../css/components/SensorInfo.less')

var SensorInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],

    getDefaultProps: function() {
        return {
            showSaveBtn: false
        };
    },
    
    getInitialState: function() {
        return {
            name: this.props.name,
            description: this.props.description,
            address: this.props.address,
            deviceId: this.props.deviceId,
            errors: this.props.errors,
            saveButtonStatus: ''
        }
    },

    toJson: function() {
        var s = this.state
        return {
            name: s.name,
            description: s.description,
            address: s.address,
            deviceId: s.deviceId,
            id: this.props.id,
            attr: this.props.attr,
        }
    },

    setErrors: function(errors) {
        this.setState({ errors: errors });
    },

    _changed: function(evt) {
        this.setState({ saveButtonStatus: '' });
        this.changed(evt, function(){
            this.props.changed && this.props.changed(this);
        }.bind(this));
    },
    
    devicePickerChanged: function(deviceId) {
        this.setState({
            deviceId: deviceId,
            saveButtonStatus: ''
        }, function() {
            this.props.changed && this.props.changed(this);
        }.bind(this));
    },

    save: function() {
        this.setState({ errors: null });
        Api.sensorUpdate(this.toJson(), function(err, sensorData) {
            if (err && !err.validation) {
                //TODO: Dispatch general error
                this.setState({ saveButtonStatus: 'error' })
                return;
            } else if (err && err.validation) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validation.errors[this.props.id]
                });
                return;
            }

            this.setState({ saveButtonStatus: 'success' });
            this.props.updatedSensor(sensorData);
        }.bind(this));
    },
    
    render: function() {
        var saveBtn;
        if (this.props.showSaveBtn && this.state.dirty) {
            saveBtn = (
                <SaveBtn
                    {...classes('save')}
                    clicked={this.save}
                    text="Save"
                    status={this.state.saveButtonStatus}/>
            );
        }

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
              <div className="pull-right">
                  {saveBtn}
              </div>
              <div style={{clear: 'both' }}></div>
            </div>
        );
    }
});
module.exports = SensorInfo;
