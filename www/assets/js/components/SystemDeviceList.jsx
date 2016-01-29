var React = require('react');
var DeviceInfo = require('./DeviceInfo.jsx');

var SystemDeviceList = React.createClass({
    getInitialState: function() {
        return {
            loading: true,
            devices: [],
            addingNew: false
        };
    },

    componentDidMount: function() {
        var self = this;
        $.ajax({
            url: '/api/v1/systems/123/devices',
            dataType: 'json',
            cache: false,
            success: function(data) {
                self.setState({devices: data, loading: false});
            },
            error: function(xhr, status, err) {
                console.error(err.toString());
            }
        });
    },

    newClicked: function() {
        //TODO: Show new device UI
    },

    render: function() {
        var deviceNodes = this.state.devices.map(function(device) {
            return (
                <DeviceInfo
                name={device.name}
                description={device.description}
                address={device.address}
                key={device.id}
                />
            );
        })
        
        var body = this.state.loading
        ? <div className="text-center"><i className="fa fa-spinner fa-spin"></i></div>
        : deviceNodes;

        return (
            <div className="cmp-DeviceList">
              <div className="header clearfix">
                <button className="btn btn-primary pull-right" onClick={this.newClicked}>New Device</button>
              </div>
              <h3 className={this.state.devices.length > 0 ? "" : " hidden"}>Devices</h3>
              {body}
            </div>
        );
    }
});
module.exports = SystemDeviceList;