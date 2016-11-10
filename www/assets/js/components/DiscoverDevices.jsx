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
        var uiFields = {}
        this.props.discoverer.uiFields.forEach(function(field) {
            uiFields[field.id] = field.default;
        });
        return {
            discovering: false,
            devices: null,
            scenes: null,
            uiFields: uiFields
        };
    },

    discover: function() {
        this.setState({
            discovering: true,
            devices: null,
            scenes: null
        });
        
        Api.discovererScanDevices(
            this.props.discoverer.id,
            this.state.uiFields,
            function(err, data) {
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
    },

    uiFieldChanged: function(id, evt) {
        var fields = Object.assign({}, this.state.uiFields);
        fields[id] = evt.target.value;
        this.setState({uiFields: fields});
    },
    
    render: function() {
        var devices
        if (this.state.devices && this.state.devices.length > 0) {
            devices = this.state.devices.map(function(device) {
                return <ImportGroup
                           key={device.id}
                           device={device}
                           createdDevice={this.props.importedDevice}
                           createdZones={this.props.importedZones}
                           createdSensors={this.props.importedSensors}
                       />;
            }.bind(this));
        }

        var importBody;
        var deviceCount = 0;
        if (this.state.devices) {
            deviceCount = this.state.devices.length;
        }

        var uiFields = this.props.discoverer.uiFields.map(function(uiField) {
            //id/name/description
            return (
                <div className="form-group" key={uiField.id}>
                    <label htmlFor={uiField.id}>{uiField.label}</label>
                    <input
                        className="form-control"
                        id={uiField.id}
                        type="text"
                        value={this.state.uiFields[uiField.id] || ""}
                        onChange={this.uiFieldChanged.bind(this, uiField.id)}
                    ></input>
                </div>
            );
        }.bind(this));
        importBody = (
            <div>
                <div {...classes('pre-import-instructions', this.props.discoverer.preScanInfo == '' ? 'hidden' : '')}>
                    {this.props.discoverer.preScanInfo}
                </div>
                <div>
                    {uiFields}
                </div>
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
        importedDevice: function(deviceJson) {
            dispatch(SystemActions.importedDevice(deviceJson));
        },
        importedZones: function(zones) {
            dispatch(ZoneActions.importedZones(zones));
        },
        importedSensors: function(sensors) {
            dispatch(SensorActions.importedSensors(sensors));
        }
    };
}

module.exports = ReactRedux.connect(null, mapDispatchToProps)(DiscoverDevices);
