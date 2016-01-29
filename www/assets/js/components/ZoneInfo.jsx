var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var SaveBtn = require('./SaveBtn.jsx');
var DevicePicker = require('./DevicePicker.jsx');
var ZoneControllerPicker = require('./ZoneControllerPicker.jsx');
var ZoneOutputPicker = require('./ZoneOutputPicker.jsx');
var ZoneTypePicker = require('./ZoneTypePicker.jsx');

module.exports = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            cid: this.getNextIdAndIncrement() + '',
            name: this.props.name,
            description: this.props.description,
            address: this.props.address,
            deviceId: this.props.deviceId,
            type: this.props.type,
            output: this.props.output,
            controller: this.props.controller,
            errors: null,
        }
    },

    toJson: function() {
        var s = this.state
        return {
            clientId: s.cid,
            name: s.name,
            description: s.description,
            address: s.address,
            deviceId: s.deviceId,
            type: s.type,
            output: s.output,
            controller: s.controller,
        }
    },

    save: function() {
        var saveBtn = this.refs.saveBtn;
        saveBtn.saving();

        this.setState({ errors: null });
        
        var self = this;
        $.ajax({
            url: '/api/v1/systems/1/zones',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(this.toJson()),
            success: function(data) {
                saveBtn.success();
            },
            error: function(xhr, status, err) {
                self.setState({ errors: (JSON.parse(xhr.responseText) || {}).errors});
                saveBtn.failure();
            }
        });            
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

    controllerChanged: function(controller) {
        this.setState({ controller: controller });
    },

    render: function() {
        return (
            <div className="cmp-ZoneInfo well">
              <div className={this.addErr('form-group', 'name')}>
                <label className="control-label" htmlFor={this.uid('name')}>Name*</label>
                <input value={this.state.name} data-statepath="name" onChange={this.changed} className="name form-control" type="text" id={this.uid('name')}/>
                {this.errMsg('name')}
              </div>
              <div className={this.addErr("form-group", 'description')}>
                <label className="control-label" htmlFor={this.uid("description")}>Description</label>
                <input value={this.state.description} data-statepath="description" onChange={this.changed} className="description form-control" type="text" id={this.uid("description")}/>
                {this.errMsg('description')}
              </div>
              <div className={this.addErr("form-group", "address")}>
                <label className="control-label" htmlFor={this.uid("address")}>Address*</label>
                <input value={this.state.address} data-statepath="address" onChange={this.changed} className="address form-control" type="text" id={this.uid("address")}/>
                {this.errMsg('address')}
              </div>
              <div className={this.addErr("form-group", "deviceId")}>
                <label className="control-label" htmlFor={this.uid("deviceId")}>Device*</label>
                <DevicePicker devices={this.props.devices} changed={this.devicePickerChanged}/>
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
              <div className={this.addErr("form-group", "controller")}>
                <label className="control-label" htmlFor={this.uid("controller")}>Controller*</label>
                <ZoneControllerPicker controller={this.props.controller} changed={this.controllerChanged}/>
                {this.errMsg('controller')}
              </div>
              <div className="clearfix">
                <button className="btn btn-primary pull-left" onClick={this.turnOn}>Turn On</button>
                <button className="btn btn-primary btnOff pull-left" onClick={this.turnOff}>Turn Off</button>
                <div className="pull-right">
                  <SaveBtn ref="saveBtn" clicked={this.save} text="Import" />
                </div>
              </div>
            </div>
        );
    }
});
