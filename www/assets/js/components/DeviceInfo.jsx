var React = require('react');
var ReactRedux = require('react-redux');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var SaveBtn = require('./SaveBtn.jsx');
var Api = require('../utils/API.js');
var ZoneInfo = require('./ZoneInfo.jsx')
var SensorInfo = require('./SensorInfo.jsx')
var Classnames = require('classnames')
var ZoneActions = require('../actions/ZoneActions.js');
var SensorActions = require('../actions/SensorActions.js');
var DeviceTypePicker = require('./DeviceTypePicker.jsx');

var DeviceInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        //TODO: need state?
        return {
            name: this.props.name || '',
            description: this.props.description || '',
            address: this.props.address,
            id: this.props.id,
            clientId: this.props.clientId,
            modelNumber: this.props.modelNumber || '',
            token: this.props.token,
            showToken: false,
            errors: null,
            saveButtonStatus: '',
            dirty: !this.props.id,
            connectionPool: this.props.connectionPool,
            cmdBuilder: this.props.cmdBuilder,
            type: this.props.type
        }
    },

    getDefaultProps: function() {
        return {
            zones: [],
            sensors: [],
            showZones: false,
            showSensors: false
        };
    },

    toJson: function() {
        var s = this.state;
        return {
            clientId: this.props.clientId,
            name: s.name,
            description: s.description,
            address: s.address,
            modelNumber: s.modelNumber,
            token: s.token,
            id: s.id,
            connPool: this.props.connectionPool,
            cmdBuilder: this.props.cmdBuilder,
            type: s.type
        };
    },

    componentWillReceiveProps: function(nextProps) {
        //TODO: Needed?
        var device = this.state.device;
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
        if (nextProps.clientId != "") {
            this.setState({ clientId: nextProps.clientId });
        }
    },

    createDevice: function() {
        Api.deviceCreate(this.toJson(), function(err, deviceData) {
            if (err) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validationErrors
                });
                return;
            }

            // Let callers know the device has been saved
            this.props.createdDevice(this.state.clientId, deviceData);
            
            // Now we need to loop through each of the zones and save them
            function saveZone(index) {
                if (index >= this.props.zones.length) {
                    saveSensor.bind(this)(0)
                    return;
                }

                // Now the device has an id, we need to bind the zone to it
                var zoneInfo = this.refs["zoneInfo_" + this.props.zones[index].clientId];
                var zone = Object.assign({}, zoneInfo.toJson());
                zone.deviceId = deviceData.id;
                Api.zoneCreate(zone, function(err, zoneData) {
                    if (err) {
                        zoneInfo.setErrors(err.validationErrors);
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

                // Now the device has an id, we need to bind the sensor to it
                var sensorInfo = this.refs["sensorInfo_" + this.props.sensors[index].clientId];
                var sensor = Object.assign({}, sensorInfo.toJson());
                sensor.deviceId = deviceData.id;
                Api.sensorCreate(sensor, function(err, sensorData) {
                    if (err) {
                        sensorInfo.setErrors(err.validationErrors);
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
    },

    updateDevice: function() {
        Api.deviceUpdate(this.toJson(), function(err, deviceData) {
            if (err) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validationErrors
                });
                return;
            }

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
        this.props.deviceDelete(this.state.id, this.state.clientId);
    },

    typeChanged: function(type) {
        //this.setState({ type: type });
        this.changed({
            target: {
                getAttribute: function() {
                    return 'type'
                },
                value: type
            }
        });
    },

    _changed: function(evt) {
        this.setState({ saveButtonStatus: '' });

        if (evt) {
            this.changed(evt);
        }
    },

    _zoneChanged: function() {
        this._changed();
    },

    _sensorChanged: function() {
        this._changed();
    },

    render: function() {
        var device = this.state.device;

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
        if (this.state.dirty) {
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
                <button className="btn btn-link btnDelete pull-right" onClick={this.deleteDevice}>
                    <i className="glyphicon glyphicon-trash"></i>
                </button>
            );
        }

        var zones;
        if (this.props.zones.length === 0) {
            zones = <h4>0 zones found</h4>
        } else {
            zones = this.props.zones.map(function(zone) {
                return (
                    <ZoneInfo
                        ref={"zoneInfo_" + zone.clientId}
                        readOnlyFields="deviceId"
                        key={zone.id || zone.clientId}
                        clientId={zone.clientId}
                        name={zone.name}
                        description={zone.description}
                        address={zone.address}
                        type={zone.type}
                        output={zone.output}
                        deviceId={this.state.id || this.state.clientId}
                        devices={[ this.toJson() ]}
                        changed={this._zoneChanged} />
                );
            }.bind(this));
        }

        var sensors;
        if (this.props.sensors.length === 0) {
            sensors = <h4>0 sensors found</h4>
        } else {
            sensors = this.props.sensors.map(function(sensor) {
                return (
                    <SensorInfo
                        ref={"sensorInfo_" + sensor.clientId}
                        readOnlyFields="deviceId"
                        key={sensor.id || sensor.clientId}
                        clientId={sensor.clientId}
                        name={sensor.name}
                        description={sensor.description}
                        address={sensor.address}
                        attr={sensor.attr}
                        deviceId={this.state.id || this.state.clientId}
                        devices={[ this.toJson() ]}
                        changed={this._sensorChanged} />
                );
            }.bind(this));
        }
            
        return (
            <div className="cmp-DeviceInfo well-sm">
                {deleteBtn}
                <div className={this.addErr("form-group", "name")}>
                    <label className="control-label" htmlFor={this.uid("name")}>Name*</label>
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
                    <label className="control-label" htmlFor={this.uid("id")}>ID</label>
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
                    <label className="control-label" htmlFor={this.uid("type")}>Type*</label>
                    <DeviceTypePicker type={this.state.type} changed={this.typeChanged}/>
                    {this.errMsg('type')}
                </div>
                <div className={this.addErr("form-group", "description")}>
                    <label className="control-label" htmlFor={this.uid("description")}>Description</label>
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
                    <label className="control-label" htmlFor={this.uid("modelNumber")}>Model Number</label>
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
                    <label className="control-label" htmlFor={this.uid("address")}>Address</label>
                    <input
                        value={this.state.address}
                        data-statepath="address"
                        onChange={this._changed}
                        className="address form-control"
                        type="text"
                        id={this.uid("address")} />
                    {this.errMsg("address")}
                </div>
                <div className={Classnames({clearfix: true, hidden: !this.props.showZones})}>
                    <a data-toggle="collapse" href={"#" + this.uid("zones")}>
                        Zones
                        <i className="glyphicon glyphicon-menu-down"></i>
                    </a>
                </div>
                <div className="collapse zones" id={this.uid("zones")}>
                    {zones}
                </div>
                <div className={Classnames({clearfix: true, hidden: !this.props.showSensors})}>
                    <a data-toggle="collapse" href={"#" + this.uid("sensors")}>
                        Sensors
                        <i className="glyphicon glyphicon-menu-down"></i>
                    </a>
                </div>
                <div className="collapse sensors" id={this.uid("sensors")}>
                    {sensors}
                </div>
                {token}
                {/*
                <button className="btn btn-primary" onClick={this.testConnection}>Test Connection</button>
                */}
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
        savedZone: function(zoneJson) {
            dispatch(ZoneActions.importedZone(zoneJson));
        },
        savedSensor: function(sensorJson) {
            dispatch(SensorActions.importedSensor(sensorJson));
        }
    };
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(DeviceInfo);
