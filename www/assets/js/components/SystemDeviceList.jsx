var React = require('react');
var DeviceInfo = require('./DeviceInfo.jsx');

var SystemDeviceList = React.createClass({
    getDefaultProps: function() {
        return {
            devices: []
        }
    },

    newClicked: function() {
        //TODO: Show new device UI
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
                  readOnlyFields="id"
                  key={device.id}
                />
            );
        })

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
module.exports = SystemDeviceList;
