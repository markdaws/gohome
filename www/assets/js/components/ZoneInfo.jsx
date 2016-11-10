var React = require('react');
var ReactRedux = require('react-redux');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var DevicePicker = require('./DevicePicker.jsx');
var ZoneOutputPicker = require('./ZoneOutputPicker.jsx');
var ZoneTypePicker = require('./ZoneTypePicker.jsx');
var SaveBtn = require('./SaveBtn.jsx');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'ZoneInfo',
    prefix: 'b-'
});
require('../../css/components/ZoneInfo.less')

//TODO: Remove individual props from this cmp, just pass in zone

var ZoneInfo = React.createClass({
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
            type: this.props.type,
            output: this.props.output,
            errors: this.props.errors,
            id: this.props.id,
            dirty: false,
            saveButtonStatus: '',
        }
    },

    toJson: function() {
        var s = this.state
        return {
            id: this.props.id,
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
        this.setState({ saveButtonStatus: '' }, function() {
            this.props.changed && this.props.changed(this);
        }.bind(this));
        this.changed(evt);
    },
    
    devicePickerChanged: function(deviceId) {
        this.setState({ deviceId: deviceId }, function() {
            this.props.changed && this.props.changed(this);
        }.bind(this));
    },

    typeChanged: function(type) {
        if (this.state.type === type) {
            return;
        }
        
        this.setState({
            saveButtonStatus: '',
            type: type,
            dirty: true
        }, function() {
            this.props.changed && this.props.changed(this);
        }.bind(this));
    },

    outputChanged: function(output) {
        if (this.state.output === output) {
            return;
        }
        
        this.setState({
            saveButtonStatus: '',
            output: output,
            dirty: true
        }, function() {
            this.props.changed && this.props.changed(this);
        }.bind(this));
    },

    save: function() {
        this.setState({ errors: null });
        Api.zoneUpdate(this.toJson(), function(err, zoneData) {
            if (err) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validationErrors
                });
                return;
            }

            this.setState({ saveButtonStatus: 'success' });
            this.props.updatedZone(zoneData);
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
              <div className={this.addErr("form-group", "type")}>
                  <label {...classes('label', '', 'control-label')} htmlFor={this.uid("type")}>Type*</label>
                  <ZoneTypePicker type={this.props.type} changed={this.typeChanged}/>
                  {this.errMsg('type')}
              </div>
              <div className={this.addErr("form-group", "output")}>
                  <label {...classes('label', '', 'control-label')} htmlFor={this.uid("output")}>Output*</label>
                  <ZoneOutputPicker output={this.props.output} changed={this.outputChanged}/>
                  {this.errMsg('output')}
              </div>
              <div className="">
                  <a data-toggle="collapse" href={"#" + this.uid("moreInfo")}>
                      More Details
                      <i className="glyphicon glyphicon-menu-down"></i>
                  </a>
              </div>
              <div className="collapse moreInfo" id={this.uid("moreInfo")}>
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
              <div className="pull-right">
                  {saveBtn}
              </div>
              <div style={{clear: 'both' }}></div>
            </div>
        );
    }
});
module.exports = ZoneInfo;
