var React = require('react');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');
var SystemDeviceListGridCell = require('./SystemDeviceListGridCell.jsx');
var ZoneSensorListGridCell = require('./ZoneSensorListGridCell.jsx');
var DeviceInfo = require('./DeviceInfo.jsx');
var Grid = require('./grid.jsx');
var ZoneInfo = require('./ZoneInfo.jsx');
var SensorInfo = require('./SensorInfo.jsx');
var SaveBtn = require('./SaveBtn.jsx');

var classes = new BEMHelper({
    name: 'ImportGroup',
    prefix: 'b-'
});
require('../../css/components/ImportGroup.less')

var ImportGroup = React.createClass({
    getInitialState: function() {
        this._deviceSaved = false;
        this._unselected = {};
        
        return {
            device: this.props.device,
            zones: this.props.device.zones,
            sensors: this.props.device.sensors,
            deviceErrors: null,
            zoneErrors: {},
            sensorErrors: {},
            saveButtonStatus: ''
        };
    },

    updateUnselected: function(id, checked) {
        if (checked) {
            delete this._unselected[id];
        } else {
            this._unselected[id] = true;
        }
    },
    
    deviceChkBxChanged: function(deviceId, checked) {
        this.updateUnselected(deviceId, checked);
    },

    zoneChkBxChanged: function(zoneId, checked) {
        this.updateUnselected(zoneId, checked);
    },

    sensorChkBxChanged: function(sensorId, checked) {
        this.updateUnselected(sensorId, checked);
    },

    _deviceChanged: function(cmp) {
        if (this._deviceSaved) {
            return;
        }
        
        this.setState({
            device: cmp.toJson(),
            saveButtonStatus: ''
        });
    },
    
    _zoneChanged: function(cmp) {
        if (this._deviceSaved) {
            return;
        }
        
        var zone = cmp.toJson();
        var zones = this.state.zones.map(function(zn) {
            if (zn.id === zone.id) {
                return zone;
            }
            return zn;
        });
        this.setState({
            zones: zones,
            saveButtonStatus: ''
        });
    },

    _sensorChanged: function(cmp) {
        if (this._deviceSaved) {
            return;
        }
        var sensor = cmp.toJson();

        // Update our list of sensors, with the newly updated sensor replacing the old version,
        // this will cause a re-render
        var sensors = this.state.sensors.map(function(sen) {
            if (sen.id === sensor.id) {
                return sensor;
            }
            return sen;
        });
        this.setState({
            sensors: sensors,
            saveButtonStatus: ''
        });
    },

    importClicked: function() {
        this.refs["devicegrid"].closeExpander();
        this.refs["zonegrid"].closeExpander();
        this.refs["sensorgrid"].closeExpander();
        this.setState({
            zoneErrors: {},
            sensorErrors: {},
            deviceErrors: null
        });

        var device = this.state.device;
        device.zones = this.state.zones.filter(function(zone) {
            return !zone.isDupe && !this._unselected[zone.id];
        }.bind(this));
        device.sensors = this.state.sensors.filter(function(sensor) {
            return !sensor.isDupe && !this._unselected[sensor.id];
        }.bind(this));

        if (device.isDupe) {
            function saveZone(index) {
                if (index >= device.zones.length) {
                    //TODO: Sensors
                    this.setState({ saveButtonStatus: 'success' });
                    return;
                }

                var zone = device.zones[index];
                Api.zoneCreate(zone, function(err, zoneData) {
                    if (err) {
                        var zoneErrs = {};
                        zoneErrs[zone.id] = err.validationErrors[zone.id];
                        this.setState({
                            saveButtonStatus: 'error',
                            zoneErrors: zoneErrs
                        });
                        return
                    }
                    this.props.createdZones([ zone ]);
                    saveZone.bind(this)(index + 1);
                }.bind(this));
            }
            saveZone.bind(this)(0);
        } else {
            Api.deviceCreate(device, function(err, deviceData) {
                if (err) {
                    var zoneErrs = {};
                    var sensorErrs = {};
                    this.state.zones.forEach(function(zone) {
                        if (err.validationErrors[zone.id]) {
                            zoneErrs[zone.id] = err.validationErrors[zone.id];
                        }
                    });
                    this.state.sensors.forEach(function(sensor) {
                        if (err.validationErrors[sensor.id]) {
                            sensorErrs[sensor.id] = err.validationErrors[sensor.id];
                        }
                    });

                    this.setState({
                        saveButtonStatus: 'error',
                        zoneErrors: zoneErrs,
                        sensorErrors: sensorErrs,
                        deviceErrors: err.validationErrors[this.state.device.id]
                    });
                    return;
                }

                this._deviceSaved = true;

                // Let callers know the device has been saved
                this.props.createdDevice(deviceData);
                this.props.createdZones(this.state.device.zones);
                this.props.createdSensors(this.state.device.sensors);

                this.setState({ saveButtonStatus: 'success' });
            }.bind(this));
        }
    },
    
    render: function() {
        var devices = [];
        var zones = [];
        var sensors = [];
        var device = this.state.device;

        // Only render non dupes
        if (!device.isDupe) {
            var cell = {
                key: device.id,
                cell: <SystemDeviceListGridCell
                          key={"devicecell-" + device.id}
                          id={device.id}
                          showCheckbox={false}
                          chkBxChanged={this.deviceChkBxChanged}
                          hasError={this.state.deviceErrors != null}
                          hasSuccess={this._deviceSaved}
                          device={device} />,
                content: <DeviceInfo
                {...device}
                key={"deviceinfo-" + device.id}
                ref={"deviceinfo-" + device.id}
                readOnlyFields="id"
                errors={this.state.deviceErrors}
                changed={this._deviceChanged} />
            };
            devices.push(cell);
        }

        (this.state.zones || []).forEach(function(zone) {
            if (zone.isDupe) {
                return;
            }
            
            // If we have an error we need to check for it
            var err = this.state.zoneErrors[zone.id];
            
            var cmpZone = {
                key: 'zones_' + zone.id,
                cell: <ZoneSensorListGridCell
                          id={zone.id}
                          showCheckbox={true}
                          showLevel={false}
                          hasError={err != null}
                          hasSuccess={this._deviceSaved && !this._unselected[zone.id]}
                          chkBxChanged={this.zoneChkBxChanged}
                          key={"zonecell-" + zone.id}
                          ref={"cell_zone_" + zone.id}
                          zone={zone} />,
                content: <ZoneInfo
                             {...zone}
                             readOnlyFields="deviceId"
                             key={"zoneinfo-" + zone.id}
                             ref={"zoneinfo-" + zone.id}
                             errors={err}
                             deviceId={device.id}
                             devices={[ device ]}
                             changed={this._zoneChanged} />
            };
            zones.push(cmpZone);
        }.bind(this));

        (this.state.sensors || []).forEach(function(sensor) {
            if (sensor.isDupe) {
                return;
            }
            
            var err = this.state.sensorErrors[sensor.id];
            
            var cmpSensor = {
                key: 'sensor_' + sensor.id,
                cell: <ZoneSensorListGridCell
                          id={sensor.id}
                          showCheckbox={true}
                          chkBxChanged={this.sensorChkBxChanged}
                          hasError={err != null}
                          hasSuccess={this._deviceSaved && !this._unselected[sensor.id]}
                          key={sensor.id}
                          ref={"cell_sensor_" + sensor.id}
                          sensor={sensor} />,
                content: <SensorInfo
                             {...sensor}
                             readOnlyFields="deviceId"
                             key={sensor.id}
                             errors={err}
                             deviceId={device.id}
                             devices={[ device ]}
                             changed={this._sensorChanged} />

            };
            sensors.push(cmpSensor);
        }.bind(this));

        return (
            <div {...classes('','', 'clearfix')}>
                <div {...classes('import', '', 'pull-right clearfix')}>
                    <SaveBtn
                        clicked={this.importClicked}
                        text="Import" status={this.state.saveButtonStatus} />
                </div>
                <div {...classes('devices', devices.length === 0 ? 'hidden' : '')}>
                    <h2 {...classes('header')}>Device</h2>
                    <Grid ref="devicegrid" key="devicegrid" cells={devices} />
                </div>
                <div key="aa" {...classes('zones', zones.length === 0 ? 'hidden' : '')}>
                    <h2 {...classes('header')}>Zones</h2>
                    <Grid ref="zonegrid" key="zonegrid" cells={zones} />
                </div>
                <div {...classes('sensors', sensors.length === 0 ? 'hidden' : '')}>
                    <h2 {...classes('header')}>Sensors</h2>
                    <Grid ref="sensorgrid" key="sensorgrid" cells={sensors} />
                </div>
                <div style={{ clear: 'both' }}></div>
            </div>
        );
    }
});
module.exports = ImportGroup;
