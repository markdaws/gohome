var React = require('react');
var ReactRedux = require('react-redux');
var DeviceInfo = require('./DeviceInfo.jsx');
var Api = require('../utils/API.js');
var SystemActions = require('../actions/SystemActions.js');

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

        Api.discoverDevice(this.props.modelNumber, function(err, data) {
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
                return <DeviceInfo
                           name={device.name}
                           description={device.description}
                           address={device.address}
                           modelNumber={device.modelNumber}
                           connectionPool={device.connPool}
                           cmdBuilder={device.cmdBuilder}
                           id={device.id}
                           clientId={device.clientId}
                           readOnlyFields="id, modelNumber"
                           key={device.id || device.clientId}
                           savedDevice={this.props.importedDevice}
                           showZones={true}
                           zones={device.zones}/>
            }.bind(this));
        }

        var importBody
        importBody = (
            <div>
                <button className={"btn btn-primary" + (this.state.discovering ? " disabled" : "")}
                        onClick={this.discover}>Discover Devices</button>
                <i className={"fa fa-spinner fa-spin discover" + (this.state.discovering ? "" : " hidden")}></i>
                <h3 className={this.state.devices ? "" : " hidden"}>
                    {this.state.devices && this.state.devices.length} device(s) found
                </h3>
                <p className={this.state.devices && this.state.devices.length > 0 ? "" : " hidden"}>
                    Click "Save" on each device you wish to add to your system.
                </p>
                {devices}
            </div>
        );
        return (
            <div className="cmp-DiscoverDevices">
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
