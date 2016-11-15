var React = require('react');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');
var SystemDeviceListGridCell = require('./SystemDeviceListGridCell.jsx');
var ZoneSensorListGridCell = require('./ZoneSensorListGridCell.jsx');
var DeviceInfo = require('./DeviceInfo.jsx');
var Grid = require('./Grid.jsx');
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
        this._unselected = {};
        this._savedDevices = {};
        this._savedZones = {};
        this._savedSensors = {};
        
        // We render all devices/zones/sensors together
        var zones = [];
        var sensors = [];
        this.props.devices.forEach(function(device) {
            device.zones.forEach(function(zone) {
                zones.push(zone);
            });
            device.sensors.forEach(function(sensor) {
                sensors.push(sensor);
            });
        });
        return {
            devices: this.props.devices,
            zones: zones,
            sensors: sensors,
            zoneErrors: {},
            sensorErrors: {},
            deviceErrors: {},
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
        var device = cmp.toJson();
        var devices = this.state.devices.map(function(dev) {
            if (dev.id === device.id) {
                return device;
            }
            return dev;
        });
        this.setState({
            devices: devices,
            saveButtonStatus: ''
        });
    },
    
    _zoneChanged: function(cmp) {
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
            deviceErrors: {},
            zoneErrors: {},
            sensorErrors: {},
        });

        // When we are saving devices, we need to sort the devices so that
        // we save devices that don't have a hub before devices that do have
        // a hub, because we can't save a device that points to a hub that we
        // haven't potentially saved yet
        var sortedDevices = ([].concat(this.state.devices)).sort(function(x, y) {
            var xHasHub = x.hubId !== '';
            var yHasHub = y.hubId !== '';

            if ((xHasHub && yHasHub) || (!xHasHub && !yHasHub)) {
                return 0;
            } else if (xHasHub && !yHasHub) {
                return 1;
            } else {
                return -1;
            }
            return x.hubId > y.hubId;
        })

        function saveDevice(devIndex) {
            if (devIndex >= sortedDevices.length) {
                this.setState({ saveButtonStatus: 'success' });
                return;
            }

            var device = sortedDevices[devIndex];

            // We grouped all of the zones/sensors together, now we need to break them
            // up in to their respective devices so they can be saved
            device.zones = this.state.zones.filter(function(zone) {
                return zone.deviceId === device.id &&
                       !zone.isDupe &&
                       !this._unselected[zone.id] &&
                       !this._savedZones[zone.id];
            }.bind(this));
            device.sensors = this.state.sensors.filter(function(sensor) {
                return sensor.deviceId === device.id &&
                       !sensor.isDupe &&
                       !this._unselected[sensor.id] &&
                       !this._savedSensors[sensor.id];
            }.bind(this));

            if (device.isDupe || this._savedDevices[device.id]) {
                // Don't save the device, save any new zones + sensors we found that aren't dupes
                function saveZone(index) {
                    if (index >= device.zones.length) {
                        // saved all of the zones, move on to the sensors
                        saveSensor.bind(this)(0);
                        return;
                    }

                    var zone = device.zones[index];
                    Api.zoneCreate(zone, function(err, zoneData) {
                        if (err) {
                            var zoneErrs = {};
                            zoneErrs[zone.id] = err.validation.errors[zone.id];
                            this.setState({
                                saveButtonStatus: 'error',
                                zoneErrors: zoneErrs
                            });
                            return
                        }
                        this.props.createdZones([ zone ]);
                        this._savedZones[zone.id] = true;
                        saveZone.bind(this)(index + 1);
                    }.bind(this));
                }
                function saveSensor(index) {
                    if (index >= device.sensors.length) {
                        // Move on to the next device
                        saveDevice.bind(this)(devIndex + 1);
                        return;
                    }

                    var sensor = device.sensors[index];
                    Api.sensorCreate(sensor, function(err, sensorData) {
                        if (err) {
                            var sensorErrs = {};
                            sensorErrs[sensor.id] = err.validation.errors[sensor.id];
                            this.setState({
                                saveButtonStatus: 'error',
                                sensorErrors: sensorErrs
                            });
                            return
                        }
                        this.props.createdSensors([ sensor ]);
                        this._savedSensors[sensor.id] = true;
                        saveSensor.bind(this)(index + 1);
                    }.bind(this));
                }
                saveZone.bind(this)(0);
            } else {
                Api.deviceCreate(device, function(err, deviceData) {
                    if (err) {
                        var zoneErrs = {};
                        var sensorErrs = {};
                        this.state.zones.forEach(function(zone) {
                            if (err.validation.errors[zone.id]) {
                                zoneErrs[zone.id] = err.validation.errors[zone.id];
                            }
                        });
                        this.state.sensors.forEach(function(sensor) {
                            if (err.validation.errors[sensor.id]) {
                                sensorErrs[sensor.id] = err.validation.errors[sensor.id];
                            }
                        });

                        this.setState({
                            saveButtonStatus: 'error',
                            zoneErrors: zoneErrs,
                            sensorErrors: sensorErrs,
                            deviceErrors: err.validation.errors[device.id] || {}
                        });
                        return;
                    }

                    // Let callers know the device has been saved
                    this.props.createdDevice(deviceData);
                    this.props.createdZones(device.zones);
                    this.props.createdSensors(device.sensors);

                    this._savedDevices[device.id] = true;
                    device.zones.forEach(function(zone) {
                        this._savedZones[zone.id] = true;
                    }.bind(this));
                    device.sensors.forEach(function(sensor) {
                        this._savedSensors[sensor.id] = true;
                    }.bind(this));

                    // Move on to the next device
                    saveDevice.bind(this)(devIndex + 1);
                }.bind(this));
            }
        }
        saveDevice.bind(this)(0);
    },

    getDeviceById: function(id) {
        var devices = this.state.devices;
        for (var i=0; i<devices.length; ++i) {
            if (devices[i].id === id) {
                return devices[i];
            }
        }
        return null;
    },
    
    render: function() {
        var devices = [];
        var zones = [];
        var sensors = [];

        (this.state.devices || []).forEach(function(device) {
            if (device.isDupe) {
                return;
            }

            var err = this.state.deviceErrors[device.id];
            var cell = {
                key: device.id,
                cell: <SystemDeviceListGridCell
                          key={"devicecell-" + device.id}
                          id={device.id}
                          showCheckbox={false}
                          chkBxChanged={this.deviceChkBxChanged}
                          hasError={err != null}
                          hasSuccess={this._savedDevices[device.id]}
                          device={device} />,
                content: <DeviceInfo
                             {...device}
                             key={"deviceinfo-" + device.id}
                             ref={"deviceinfo-" + device.id}
                             readOnlyFields="id"
                             errors={err}
                             changed={this._deviceChanged} />
            };
            devices.push(cell);
        }.bind(this));

        (this.state.zones || []).forEach(function(zone) {
            if (zone.isDupe) {
                return;
            }
            
            var err = this.state.zoneErrors[zone.id];
            var device = this.getDeviceById(zone.deviceId);
            var cmpZone = {
                key: 'zones_' + zone.id,
                cell: <ZoneSensorListGridCell
                          id={zone.id}
                          showCheckbox={true}
                          showLevel={false}
                          hasError={err != null}
                          hasSuccess={this._savedZones[zone.id]}
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
            var device = this.getDeviceById(sensor.deviceId);
            var cmpSensor = {
                key: 'sensor_' + sensor.id,
                cell: <ZoneSensorListGridCell
                          id={sensor.id}
                          showCheckbox={true}
                          chkBxChanged={this.sensorChkBxChanged}
                          hasError={err != null}
                          hasSuccess={this._savedSensors[sensor.id]}
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

        var hasNonDupes = devices.length > 0 || zones.length > 0 || sensors.length > 0;
        var body;
        if (hasNonDupes) {
            body = (
                <div>
                    <div {...classes('import', '', 'pull-right clearfix')}>
                        <SaveBtn
                            clicked={this.importClicked}
                            text="Import" status={this.state.saveButtonStatus} />
                    </div>
                    <div {...classes('devices', devices.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('header')}>Devices</h2>
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
        } else {
            body = (
                <div {...classes('no-new')}>
                    No new devices/zones/sensors found. Items previously imported will not be shown again
                    unless you delete them from the system
                </div>
            );
        }
        return (
            <div {...classes('','', 'clearfix')}>
                {body}
            </div>
        );
    }
});
module.exports = ImportGroup;
