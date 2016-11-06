var React = require('react');
var ReactRedux = require('react-redux');
var DeviceInfo = require('./DeviceInfo.jsx');
var Api = require('../utils/API.js');
var SystemActions = require('../actions/SystemActions.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'DiscoverDevices',
    prefix: 'b-'
});
require('../../css/components/DiscoverDevices.less')

var DiscoverDevices = React.createClass({
    getInitialState: function() {
        return {
            discovering: false,
            devices: null,
        };
    },

    discover: function() {
        this.setState({
            discovering: true,
            devices: null
        });

        Api.discovererScanDevices(this.props.discoverer.id, function(err, data) {
            this.setState({
                discovering: false,
                devices: data || []
            });
        }.bind(this));
    },

    render: function() {
        var devices
        if (this.state.devices && this.state.devices.length > 0) {
            devices = this.state.devices.map(function(device) {
                return (
                    <div {...classes('device-info')} key={device.id || device.clientId}>
                        <DeviceInfo
                           name={device.name}
                           description={device.description}
                           address={device.address}
                           modelNumber={device.modelNumber}
                           modelName={device.modelName}
                           softwareVersion={device.softwareVersion}
                           connectionPool={device.connPool}
                           cmdBuilder={device.cmdBuilder}
                           auth={device.auth}
                           id={device.id}
                           clientId={device.clientId}
                           readOnlyFields="id, modelNumber"
                           key={device.id || device.clientId}
                           createdDevice={this.props.importedDevice}
                           showZones={true}
                           showSensors={true}
                           zones={device.zones}
                           sensors={device.sensors}/>
                    </div>
                );
            }.bind(this));
        }

        var importBody
        importBody = (
            <div>
                <div {...classes('pre-import-instructions', this.props.discoverer.preScanInfo == '' ? 'hidden' : '')}>
                    {this.props.discoverer.preScanInfo}
                </div>
                <div {...classes('discover')}>
                    <button {...classes('', '', (this.state.discovering ? 'disabled' : '') + ' btn btn-primary')}
                        onClick={this.discover}>Discover Devices</button>
                    <i {...classes('spinner', this.state.discovering ? '' : 'hidden', 'fa fa-spinner fa-spin')}></i>
                </div>
                <h3 {...classes('no-devices', this.state.devices ? '' : 'hidden')}>
                    {this.state.devices && this.state.devices.length} device(s) found
                </h3>
                <p {...classes('found-devices', this.state.devices && this.state.devices.length > 0 ? '' : ' hidden')}>
                    Click "Save" on each device you wish to add to your system.
                </p>
                {devices}
            </div>
        );
        return (
            <div {...classes()}>
                {importBody}
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        importedDevice: function(clientId, deviceJson) {
            dispatch(SystemActions.importedDevice(deviceJson));
        }
    };
}

module.exports = ReactRedux.connect(null, mapDispatchToProps)(DiscoverDevices);
