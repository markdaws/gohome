var React = require('react');
var ReactRedux = require('react-redux');
var DeviceInfo = require('./DeviceInfo.jsx');
var Api = require('../utils/API.js');
var SystemActions = require('../actions/SystemActions.js');
var ZoneActions = require('../actions/ZoneActions.js');
var SensorActions = require('../actions/SensorActions.js');
var ImportGroup = require('./ImportGroup.jsx');
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
            scenes: null,
            configString: ''
        };
    },

    discover: function() {
        var configString = '';
        try {
            // Remove all the whitespace, otherwise go gets upset
            configString = JSON.stringify(JSON.parse(this.state.configString))
        }
        catch(e) {
            console.error('TODO: show error');
            this.setState();
        }
        this.setState({
            discovering: true,
            devices: null,
            scenes: null
        });

        if (this.props.discoverer.type === 'FromString') {
            Api.discovererFromString(this.props.discoverer.id, configString, function(err, data) {
                if (err != null) {
                    //TODO: Change import button to green/red
                    console.error("Failed to discover");
                    console.error(err);
                    return;
                }

                this.setState({
                    discovering: false,
                    scenes: data.scenes,
                    devices: data.devices
                });            
            }.bind(this));
        } else {
            Api.discovererScanDevices(this.props.discoverer.id, function(err, data) {
                if (err != null) {
                    //TODO: Change import button to green/red
                    console.error("Failed to discover");
                    console.error(err);
                    return;
                }
                this.setState({
                    discovering: false,
                    scenes: data.scenes,
                    devices: data.devices
                });
            }.bind(this));
        }
    },

    textChanged: function(evt) {
        this.setState({ configString: evt.target.value });
    },
    
    render: function() {
        var devices
        if (this.state.devices && this.state.devices.length > 0) {
            devices = this.state.devices.map(function(device) {
                return <ImportGroup
                           key={device.id || device.clientId}
                           device={device}
                           createdDevice={this.props.importedDevice}
                           createdZone={this.props.importedZone}
                           createdSensor={this.props.importedSensor}
                       />;
            }.bind(this));
        }

        var importBody;
        var deviceCount = 0;
        if (this.state.devices) {
            deviceCount = this.state.devices.length;
        }
        
        importBody = (
            <div>
                <div {...classes('pre-import-instructions', this.props.discoverer.preScanInfo == '' ? 'hidden' : '')}>
                    {this.props.discoverer.preScanInfo}
                </div>
                <textarea
                    {...classes('text-area', this.props.discoverer.type === 'FromString' ? '' : 'hidden')}
                    placeholder="Paste the config string here" value={this.state.configString} onChange={this.textChanged}>
                </textarea>
                <div {...classes('discover')}>
                    <button {...classes('', '', (this.state.discovering ? 'disabled' : '') + ' btn btn-primary')}
                        onClick={this.discover}>Discover Devices</button>
                    <i {...classes('spinner', this.state.discovering ? '' : 'hidden', 'fa fa-spinner fa-spin')}></i>
                </div>
                <h3 {...classes('no-devices', this.state.devices && deviceCount === 0 ? '' : 'hidden')}>
                    {deviceCount} device{deviceCount > 1 ? 's' : ''} found
                </h3>
                <p {...classes('found-devices', this.state.devices && this.state.devices.length > 0 ? '' : ' hidden')}>
                    Click "Import" on each device you wish to add to your system. Uncheck the check boxes next to items
                    you do not want to import.
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
        },
        importedZone: function(zoneJson) {
            dispatch(ZoneActions.importedZone(zoneJson));
        },
        importedSensor: function(sensorJson) {
            dispatch(SensorActions.importedSensor(sensorJson));
        }
    };
}

module.exports = ReactRedux.connect(null, mapDispatchToProps)(DiscoverDevices);
