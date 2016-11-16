var React = require('react');
var ReactRedux = require('react-redux');
var DeviceInfo = require('./DeviceInfo.jsx');
var SystemActions = require('../actions/SystemActions.js');
var Grid = require('./Grid.jsx');
var DeviceCell = require('./DeviceCell.jsx')
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'DeviceList',
    prefix: 'b-'
});
require('../../css/components/DeviceList.less')

var DeviceList = React.createClass({
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
                cell: <DeviceCell
                          key={device.id}
                          device={device} />,
                content: <DeviceInfo
                             name={device.name}
                             key={device.id}
                             showSaveBtn={true}
                             description={device.description}
                             address={device.address}
                             modelNumber={device.modelNumber}
                             modelName={device.modelName}
                             softwareVersion={device.softwareVersion}
                             id={device.id}
                             readOnlyFields="id"
                             key={device.id}
                             type={device.type}
                             deviceDelete={this.props.deviceDelete}
                             createdDevice={this.props.createdDevice}
                             updatedDevice={this.props.updatedDevice}/>
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

        var dimmerSection;
        if (dimmers.length > 0) {
            dimmerSection = (
                <div key="dimmerSection">
                    <h2 {...classes('header')}>Dimmers</h2>
                    <Grid cells={dimmers} />
                </div>
            );
        }
        var switchSection;
        if (switches.length > 0) {
            switchSection = (
                <div key="switchSection">
                    <h2 {...classes('header')}>Switches</h2>
                    <Grid cells={switches} />
                </div>
            );
        }
        var shadeSection;
        if (shades.length > 0) {
            shadeSection = (
                <div key="shadeSection">
                    <h2 {...classes('header')}>Shades</h2>
                    <Grid cells={shades} />
                </div>
            );
        }
        var hubSection;
        if (hubs.length > 0) {
            hubSection = (
                <div className={hubs.length > 0 ? "" : " hidden"}>
                    <h2 {...classes('header')}>Hubs</h2>
                    <Grid cells={hubs} />
                </div>
            );
        }
        var remoteSection;
        if (remotes.length > 0) {
            remoteSection = (
                <div key="remoteSection">
                    <h2 {...classes('header')}>Remotes</h2>
                    <Grid cells={remotes} />
                </div>
            );
        }
        var deviceSection;
        if (unknown.length > 0) {
            deviceSection = (
                <div key="deviceSection">
                    <h2 {...classes('header')}>Devices</h2>
                    <Grid cells={unknown} />
                </div>
            );
        }
        return (
            <div {...classes()}>
                {dimmerSection}
                {switchSection}
                {shadeSection}
                {hubSection}
                {remoteSection}
                {deviceSection}
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        deviceDelete: function(id) {
            dispatch(SystemActions.deviceDelete(id));
        },
        createdDevice: function(id, data) {
            dispatch(SystemActions.createdDevice(id, data));
        },
        updatedDevice: function(data) {
            dispatch(SystemActions.updatedDevice(data));
        }
    };
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(DeviceList);
