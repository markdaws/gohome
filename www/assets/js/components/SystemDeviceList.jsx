var React = require('react');
var ReactRedux = require('react-redux');
var DeviceInfo = require('./DeviceInfo.jsx');
var SystemActions = require('../actions/SystemActions.js');
var Grid = require('./Grid.jsx');
var SystemDeviceListGridCell = require('./SystemDeviceListGridCell.jsx')

var SystemDeviceList = React.createClass({
    getDefaultProps: function() {
        return {
            devices: []
        }
    },

    render: function() {
        var switches = [], shades = [], dimmers = [], hubs = [], remotes = [], unknown = [];
        this.props.devices.forEach(function(device) {
            var cell = {
                key: device.id,
                cell: <SystemDeviceListGridCell
                          key={device.id}
                          device={device} />,
                content: <DeviceInfo
                             name={device.name}
                             key={device.id}
                             description={device.description}
                             address={device.address}
                             modelNumber={device.modelNumber}
                             id={device.id}
                             clientId={device.clientId}
                             readOnlyFields="id"
                             key={device.id || device.clientId}
                             type={device.type}
                             deviceDelete={this.props.deviceDelete}
                             savedDevice={this.props.savedDevice} />
            };

            switch (device.type) {
                case 'dimmer':
                    dimmers.push(cell)
                    break;
                case 'switch':
                    switches.push(cell)
                    break;
                case 'shade':
                    shades.push(cell);
                    break;
                case 'hub':
                    hubs.push(cell)
                    break;
                case 'remote':
                    remotes.push(cell)
                    break;
                default:
                    unknown.push(cell)
                    break;
            }
        }.bind(this));

        return (
            <div className="cmp-SystemDeviceList">
                <div className={dimmers.length > 0 ? "" : " hidden"}>
                    <h2>Dimmers</h2>
                    <Grid cells={dimmers} />
                </div>

                <div className={switches.length > 0 ? "" : " hidden"}>
                    <h2>Switches</h2>
                    <Grid cells={switches} />
                </div>
                <div className={shades.length > 0 ? "" : " hidden"}>
                    <h2>Shades</h2>
                    <Grid cells={shades} />
                </div>
                <div className={hubs.length > 0 ? "" : " hidden"}>
                    <h2>Hubs</h2>
                    <Grid cells={hubs} />
                </div>
                <div className={remotes.length > 0 ? "" : " hidden"}>
                    <h2>Remotes</h2>
                    <Grid cells={remotes} />
                </div>
                <div className={unknown.length > 0 ? "" : " hidden"}>
                    <h2>Devices</h2>
                    <Grid cells={unknown} />
                </div>
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        deviceDelete: function(id, clientId) {
            dispatch(SystemActions.deviceDelete(id, clientId));
        },
        savedDevice: function(clientId, data) {
            dispatch(SystemActions.savedDevice(clientId, data));
        }
    };
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(SystemDeviceList);
