var React = require('react');
var ReactRedux = require('react-redux');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var SaveBtn = require('./SaveBtn.jsx');
var Api = require('../utils/API.js');
var DeviceTypePicker = require('./DeviceTypePicker.jsx');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'DeviceInfo',
    prefix: 'b-'
});
require('../../css/components/DeviceInfo.less')

var DeviceInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        //TODO: need state?
        return {
            name: this.props.name || '',
            description: this.props.description || '',
            address: this.props.address,
            id: this.props.id,
            hubId: this.props.hubId,
            modelNumber: this.props.modelNumber || '',
            modelName: this.props.modelName || '',
            softwareVersion: this.props.softwareVersion || '',
            auth: this.props.auth,
            showToken: false,
            errors: this.props.errors,
            saveButtonStatus: '',
            dirty: !this.props.id,
            connPool: this.props.connPool,
            cmdBuilder: this.props.cmdBuilder,
            type: this.props.type,
        }
    },

    getDefaultProps: function() {
        return {
            showSaveBtn: false,
        };
    },

    toJson: function() {
        var s = this.state;
        return {
            id: s.id,
            name: s.name,
            description: s.description,
            address: s.address,
            modelNumber: s.modelNumber,
            modelName: s.modelName,
            softwareVersion: s.softwareVersion,
            auth: s.auth,
            hubId: s.hubId,
            buttons: this.props.buttons,
            connPool: this.props.connPool,
            cmdBuilder: this.props.cmdBuilder,
            type: s.type
        };
    },

    componentWillReceiveProps: function(nextProps) {
        //TODO: Needed?
        if (nextProps.name != "") {
            this.setState({ name: nextProps.name });
        }
        if (nextProps.description != "") {
            this.setState({ description: nextProps.description });
        }
        if (nextProps.address != "") {
            this.setState({ address: nextProps.address });
        }
        if (nextProps.type != "") {
            this.setState({ type: nextProps.type });
        }
        if (nextProps.token != "") {
            this.setState({ token: nextProps.token });
        }
        if (nextProps.id != "" ) {
            this.setState({ id: nextProps.id });
        }
    },

    createDevice: function() {
        /*
        //TODO: Revisit now ew have one API to save everything at once
        Api.deviceCreate(this.toJson(), function(err, deviceData) {
            if (err) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validation.errors
                });
                return;
            }

            // Let callers know the device has been saved
            this.props.createdDevice(this.state.id, deviceData);

            // Now we need to loop through each of the zones and save them
            function saveZone(index) {
                if (index >= this.props.zones.length) {
                    saveSensor.bind(this)(0)
                    return;
                }

                var zone = this.refs["zoneInfo_" + this.props.zones[index].id].toJson();
                Api.zoneCreate(zone, function(err, zoneData) {
                    if (err) {
                        zoneInfo.setErrors(err.validation.errors);
                        this.setState({
                            saveButtonStatus: 'error'
                        });
                        return;
                    }

                    this.props.savedZone(zoneData);
                    saveZone.bind(this)(index+1);
                }.bind(this));
            }
            saveZone.bind(this)(0);

            // Loop through sensors saving
            function saveSensor(index) {
                if (index >= this.props.sensors.length) {
                    this.setState({ saveButtonStatus: 'success' });
                    return;
                }

                var sensor = this.refs["sensorInfo_" + this.props.sensors[index].id].toJson();
                Api.sensorCreate(sensor, function(err, sensorData) {
                    if (err) {
                        sensorInfo.setErrors(err.validation.errors);
                        this.setState({
                            saveButtonStatus: 'error'
                        });
                        return;
                    }

                    this.props.savedSensor(sensorData);
                    saveSensor.bind(this)(index+1);
                }.bind(this));
            }

        }.bind(this));
        */
    },

    updateDevice: function() {
        Api.deviceUpdate(this.toJson(), function(err, deviceData) {
            if (err && !err.validation) {
                //TODO: Dispatch general error so it can be displayed somewhere in the UI ...
                this.setState({
                    saveButtonStatus: 'error',
                });
                return;
            } else if (err && err.validation) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validation.errors[this.state.id]
                });
                return;
            }

            this.setState({ saveButtonStatus: 'success' });
            this.props.updatedDevice(deviceData);
        }.bind(this));
    },

    save: function() {
        this.setState({ errors: null });

        if (this.state.id) {
            this.updateDevice();
        } else {
            this.createDevice();
        }
    },

    deleteDevice: function() {
        this.props.deviceDelete(this.state.id);
    },

    typeChanged: function(type) {
        this.changed({
            target: {
                getAttribute: function() {
                    return 'type'
                },
                value: type
            }
        }, function() {
            this.props.changed && this.props.changed(this);
        });
    },

    _changed: function(evt) {
        this.setState({ saveButtonStatus: '' });

        if (evt) {
            this.changed(evt, function() {
                this.props.changed && this.props.changed(this)
            }.bind(this));
        }
    },

    _zoneChanged: function() {
        this._changed();
    },

    _sensorChanged: function() {
        this._changed();
    },

    render: function() {
        var token
        if (this.props.showToken) {
            token = (
                <div className={this.addErr("form-group", "token")}>
                    <label className="control-label" htmlFor={this.uid("token")}>Security Token</label>
                    <input
                        value={this.state.token}
                        data-statepath="token"
                        onChange={this._changed}
                        className="token form-control"
                        type="text"
                        id={this.uid("token")} />
                    {this.errMsg('token')}
                </div>
            );
        }

        var saveBtn;
        if (this.props.showSaveBtn && this.state.dirty) {
            saveBtn = (
                <SaveBtn
                    clicked={this.save}
                    text="Save"
                    status={this.state.saveButtonStatus}/>
            );
        }

        var deleteBtn
        if (this.props.deleteDevice) {
            deleteBtn = (
                <button {...classes('delete', '', 'btn btn-link pull-right')} onClick={this.deleteDevice}>
                    <i className="glyphicon glyphicon-trash"></i>
                </button>
            );
        }

        return (
            <div {...classes('', '', 'well well-sm')}>
                {deleteBtn}
                <div className={this.addErr("form-group", "name")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("name")}>Name*</label>
                    <input
                        value={this.state.name}
                        data-statepath="name"
                        onChange={this._changed}
                        className="name form-control"
                        type="text"
                        id={this.uid("name")} />
                    {this.errMsg("name")}
                </div>
                <div className={this.addErr("form-group", "id")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("id")}>ID</label>
                    <input
                        value={this.state.id}
                        readOnly={this.isReadOnly("id")}
                        data-statepath="id"
                        onChange={this._changed}
                        className="id form-control"
                        type="text"
                        id={this.uid("id")} />
                    {this.errMsg("id")}
                </div>
                <div className={this.addErr("form-group", "type")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("type")}>Type*</label>
                    <DeviceTypePicker type={this.state.type} changed={this.typeChanged}/>
                    {this.errMsg('type')}
                </div>
                <div className={this.addErr("form-group", "description")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("description")}>Description</label>
                    <input
                        value={this.state.description}
                        data-statepath="description"
                        onChange={this._changed}
                        className="description form-control"
                        type="text"
                        id={this.uid("description")} />
                    {this.errMsg("description")}
                </div>
                <div className={this.addErr("form-group", "modelNumber")}>
                    <label {...classes('label', '', 'control-label')} htmlFor={this.uid("modelNumber")}>Model Number</label>
                    <input
                        value={this.state.modelNumber}
                        readOnly={this.isReadOnly("modelNumber")}
                        data-statepath="modelNumber"
                        onChange={this._changed}
                        className="modelNumber form-control"
                        type="text"
                        id={this.uid("modelNumber")} />
                    {this.errMsg("modelNumber")}
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
                    {this.errMsg("address")}
                </div>
                {token}
                <div className="pull-right">
                    {saveBtn}
                </div>
                <div style={{clear:"both"}}></div>
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
    };
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(DeviceInfo);
