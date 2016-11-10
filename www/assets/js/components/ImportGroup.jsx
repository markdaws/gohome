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
        this._savedZones = {};
        this._savedSensors = {};
        this._deviceSaved = false;
        this._deviceId = null;
        
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

    deviceChkBxChanged: function(checked) { },

    zoneChkBxChanged: function(checked) { },

    sensorChkBxChanged: function(checked) { },

    _deviceChanged: function(cmp) {
        this.setState({
            device: cmp.toJson(),
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
            if (sen.id === sensor.idd) {
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

        // Now we need to loop through each of the zones and save them
        function saveZone(index) {
            if (index >= this.state.zones.length) {
                saveSensor.bind(this)(0)
                return;
            }

            var zone = this.state.zones[index];
            var zoneCell = this.refs["cell_zone_" + zone.id];
            if (!zoneCell.isChecked()) {
                saveZone.bind(this)(index + 1);
                return;
            }

            if (this._savedZones[zone.id]) {
                // Already saved this, on a previous call before we ran in to
                // an error, skip and keep going
                saveZone.bind(this)(index + 1);
                return;
            }
            
            Api.zoneCreate(zone, function(err, zoneData) {
                if (err) {
                    var errs = {};
                    errs[zone.id] = err.validationErrors;
                    this.setState({
                        zoneErrors: errs,
                        saveButtonStatus: 'error'
                    });
                    return;
                }

                this._savedZones[zone.id] = true;
                this.props.createdZone(zoneData);
                saveZone.bind(this)(index+1);
            }.bind(this));
        }

        // Loop through sensors saving
        function saveSensor(index) {
            if (index >= this.state.sensors.length) {
                this.setState({ saveButtonStatus: 'success' });
                return;
            }

            var sensor = this.state.sensors[index];
            var sensorCell = this.refs["cell_sensor_" + sensor.id];
            if (!sensorCell.isChecked()) {
                saveSensor.bind(this)(index + 1);
                return;
            }

            if (this._savedSensors[sensor.id]) {
                saveSensor.bind(this)(index + 1);
                return;
            }

            Api.sensorCreate(sensor, function(err, sensorData) {
                if (err) {
                    var errs = {};
                    errs[sensor.id] = err.validationErrors;
                    this.setState({
                        sensorErrors: errs,
                        saveButtonStatus: 'error'
                    });
                    return;
                }

                this._savedSensors[sensor.id] = true;
                this.props.createdSensor(sensorData);
                saveSensor.bind(this)(index+1);
            }.bind(this));
        }

        if (this._deviceSaved) {
            // Start saving zones and sensors
            saveZone.bind(this)(0);
        } else {
            Api.deviceCreate(this.state.device, function(err, deviceData) {
                if (err) {
                    this.setState({
                        saveButtonStatus: 'error',
                        deviceErrors: err.validationErrors
                    });
                    return;
                }

                this._deviceSaved = true;

                // Let callers know the device has been saved
                this.props.createdDevice(deviceData);

                // Now save the zones and sensors
                saveZone.bind(this)(0);
            }.bind(this));
        }
    },
    
    render: function() {
        var devices = [];
        var zones = [];
        var sensors = [];

        var device = this.state.device;
        var cell = {
            key: device.id,
            cell: <SystemDeviceListGridCell
                      key={"devicecell-" + device.id}
                      showCheckbox={false}
                      chkBxChanged={this.deviceChkBxChanged}
                      hasError={this.state.deviceErrors != null}
                      device={device} />,
            content: <DeviceInfo
                         name={device.name}
                         key={"deviceinfo-" + device.id}
                         ref={"deviceinfo-" + device.id}
                         description={device.description}
                         address={device.address}
                         addressRequired={device.addressRequired}
                         modelNumber={device.modelNumber}
                         modelName={device.modelName}
                         softwareVersion={device.softwareVersion}
                         id={device.id}
                         readOnlyFields="id"
                         errors={this.state.deviceErrors}
                         changed={this._deviceChanged}
                         type={device.type} />
        };
        devices.push(cell);

        (this.state.zones || []).forEach(function(zone) {
            // If we have an error we need to check for it
            var err = this.state.zoneErrors[zone.id];
            
            var cmpZone = {
                key: 'zones_' + zone.id,
                cell: <ZoneSensorListGridCell
                          showCheckbox={true}
                          showLevel={false}
                          hasError={err != null}
                          chkBxChanged={this.zoneChkBxChanged}
                          key={"zonecell-" + zone.id}
                          ref={"cell_zone_" + zone.id}
                          zone={zone} />,
                content: <ZoneInfo
                             readOnlyFields="deviceId"
                             key={"zoneinfo-" + zone.id}
                             ref={"zoneinfo-" + zone.id}
                             name={zone.name}
                             description={zone.description}
                             address={zone.address}
                             type={zone.type}
                             output={zone.output}
                             errors={err}
                             deviceId={device.id}
                             devices={[ device ]}
                             changed={this._zoneChanged} />
            };
            zones.push(cmpZone);
        }.bind(this));

        (this.state.sensors || []).forEach(function(sensor) {
            var err = this.state.sensorErrors[sensor.id];
            
            var cmpSensor = {
                key: 'sensor_' + sensor.id,
                cell: <ZoneSensorListGridCell
                          showCheckbox={true}
                          chkBxChanged={this.sensorChkBxChanged}
                          hasError={err != null}
                          key={sensor.id}
                          ref={"cell_sensor_" + sensor.id}
                          sensor={sensor} />,
                content: <SensorInfo
                             readOnlyFields="deviceId"
                             key={sensor.id}
                             name={sensor.name}
                             description={sensor.description}
                             address={sensor.address}
                             attr={sensor.attr}
                             errors={err}
                             deviceId={device.id}
                             devices={[ device ]}
                             changed={this._sensorChanged} />

            };
            sensors.push(cmpSensor);
        }.bind(this));

        return (
            <div {...classes()}>
                <div {...classes('import', '', 'pull-right clearfix')}>
                    <SaveBtn
                        className="whatcha"
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
            </div>
        );
    }
});
module.exports = ImportGroup;
