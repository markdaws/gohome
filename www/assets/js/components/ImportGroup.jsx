var React = require('react');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');
var SystemDeviceListGridCell = require('./SystemDeviceListGridCell.jsx');
var ZoneSensorListGridCell = require('./ZoneSensorListGridCell.jsx');
var DeviceInfo = require('./DeviceInfo.jsx');
var Grid = require('./grid.jsx');
var ZoneInfo = require('./ZoneInfo.jsx');
var SensorInfo = require('./SensorInfo.jsx');

var classes = new BEMHelper({
    name: 'ImportGroup',
    prefix: 'b-'
});
require('../../css/components/ImportGroup.less')

var ImportGroup = React.createClass({
    getInitialState: function() {
        return null;
    },

    deviceChkBxChanged: function(checked) {
        console.log('i changed - device');
    },

    zoneChkBxChanged: function(checked) {
        console.log('i changed - zone');
    },

    sensorChkBxChanged: function(checked) {
        console.log('i changed - sensor');
    },

    _zoneChanged: function() {
        //TODO:
    },

    _sensorChanged: function() {
        //TODO:
    },

    render: function() {
        var devices = [];
        var zones = [];
        var sensors = [];

        var device = this.props.device;
        var cell = {
            key: device.id,
            cell: <SystemDeviceListGridCell
                      key={device.id}
                      showCheckbox={true}
                      chkBxChanged={this.deviceChkBxChanged}
                      device={device} />,
            content: <DeviceInfo
                         name={device.name}
                         key={device.id}
                         description={device.description}
                         address={device.address}
                         modelNumber={device.modelNumber}
                         modelName={device.modelName}
                         softwareVersion={device.softwareVersion}
                         id={device.id}
                         clientId={device.clientId}
                         readOnlyFields="id"
                         key={device.id || device.clientId}
                         type={device.type}
                         deviceDelete={this.props.deviceDelete}
                         createdDevice={this.props.createdDevice}
                         updatedDevice={this.props.updatedDevice}/>
        };
        devices.push(cell);

        (device.zones || []).forEach(function(zone) {
            var cmpZone = {
                key: 'zones_' + zone.id,
                cell: <ZoneSensorListGridCell
                          showCheckbox={true}
                          chkBxChanged={this.zoneChkBxChanged}
                          key={zone.id}
                          ref={"cell_zone_" + zone.clientId}
                          zone={zone} />,
                content: <ZoneInfo
                             ref={"zoneInfo_" + zone.clientId}
                             readOnlyFields="deviceId"
                             key={zone.clientId}
                             name={zone.name}
                             description={zone.description}
                             address={zone.address}
                             type={zone.type}
                             output={zone.output}
                             deviceId={device.clientId}
                             devices={[ device ]}
                             changed={this._zoneChanged} />
            };
            zones.push(cmpZone);
        }.bind(this));

        (device.sensors || []).forEach(function(sensor) {
            var cmpSensor = {
                key: 'sensor_' + sensor.clientId,
                cell: <ZoneSensorListGridCell
                          showCheckbox={true}
                          chkBxChanged={this.sensorChkBxChanged}
                          key={sensor.id}
                          ref={"cell_sensor_" + sensor.id}
                          sensor={sensor} />,
                content: <SensorInfo
                             ref={"sensorInfo_" + sensor.clientId}
                             readOnlyFields="deviceId"
                             key={sensor.clientId}
                             name={sensor.name}
                             description={sensor.description}
                             address={sensor.address}
                             attr={sensor.attr}
                             deviceId={device.clientId}
                             devices={[ device ]}
                             changed={this._sensorChanged} />

            };
            sensors.push(cmpSensor);
        }.bind(this));

        return (
            <div {...classes()}>
                <button {...classes('import', '', 'btn btn-primary pull-right')}>Import</button>
                <div {...classes('devices', devices.length === 0 ? 'hidden' : '')}>
                    <h2 {...classes('header')}>Device</h2>
                    <Grid cells={devices} />
                </div>
                <div {...classes('zones', zones.length === 0 ? 'hidden' : '')}>
                    <h2 {...classes('header')}>Zones</h2>
                    <Grid cells={zones} />
                </div>
                <div {...classes('sensors', sensors.length === 0 ? 'hidden' : '')}>
                    <h2 {...classes('header')}>Sensors</h2>
                    <Grid cells={sensors} />
                </div>
            </div>
        );
    }
});
module.exports = ImportGroup;
