var React = require('react');
var ReactRedux = require('react-redux');
var DeviceInfo = require('./DeviceInfo.jsx');
var SystemActions = require('../actions/SystemActions.js');

var SystemDeviceList = React.createClass({
    getDefaultProps: function() {
        return {
            devices: []
        }
    },

    newClicked: function() {
        this.props.deviceNew();
    },

    render: function() {
        var deviceNodes = this.props.devices.map(function(device) {
            return (
                <DeviceInfo
                    name={device.name}
                    description={device.description}
                    address={device.address}
                    modelNumber={device.modelNumber}
                    id={device.id}
                    clientId={device.clientId}
                    readOnlyFields="id"
                    key={device.id || device.clientId}
                    deviceDelete={this.props.deviceDelete}/>
            );
        }.bind(this));

        return (
            <div className="cmp-DeviceList">
                <div className="header clearfix">
                    <button className="btn btn-primary pull-right" onClick={this.newClicked}>New Device</button>
                </div>
                <h3 className={this.props.devices.length > 0 ? "" : " hidden"}>Devices</h3>
                {deviceNodes}
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        deviceNew: function() {
            dispatch(SystemActions.deviceNew());
        },
        deviceDelete: function(id, clientId) {
            dispatch(SystemActions.deviceDelete(id, clientId));
        },
        savedDevice: function(data) {
            dispatch(SystemActions.savedDevice());
        }
    };
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(SystemDeviceList);
