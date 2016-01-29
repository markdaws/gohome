var React = require('react');
var ZoneInfo = require('./ZoneInfo.jsx');

var ImportFluxWIFI = React.createClass({
    getInitialState: function() {
        return {
            discovering: false,
            zones: [],
            loading: true,
            devices: [],
        };
    },

    componentDidMount: function() {
        var self = this;
        $.ajax({
            url: '/api/v1/systems/123/devices',
            dataType: 'json',
            cache: false,
            success: function(data) {
                self.filterDevices(data || []);
            },
            error: function(xhr, status, err) {
                console.error(err.toString());
            }
        });
    },

    filterDevices: function(devices) {
        var filteredDevices = [];
        for (var i=0; i<devices.length; ++i) {
            switch(devices[i].modelNumber) {
            default:
                //TODO: undo
                //                case 'GoHomeHub':
                filteredDevices.push(devices[i]);
                break;
            }
        }

        this.setState({
            devices: filteredDevices,
            loading: false
        });
    },
    
    discover: function() {
        this.setState({
            discovering: true,
            zones: []
        });

        var self = this;
        $.ajax({
            url: '/api/v1/discovery/FluxWIFI/zones',
            dataType: 'json',
            cache: false,
            success: function(data) {
                self.setState({
                    discovering: false,
                    zones: data
                });
            },
            error: function(xhr, status, err) {
                self.setState({
                    discovering: false
                });
                console.error(err);
            }
        });
    },

    render: function() {

        var loading
        if (this.state.loading) {
            loading = <div className="spinnerWrapper">
            <i className="fa fa-spinner fa-spin"></i></div>
        }

        var noDeviceBody
        if (!this.state.loading && this.state.devices.length === 0) {
            noDeviceBody = (
                <div>
                  <h3>Import failed</h3>
                  <p>In order to import Flux WIFI bulbs, you must have a device in your system
                  that is capable of controlling them.  Please add one of the following devices
                  to your system first, then come back and try to import again:
                  </p>
                  <ul>
                    <li>GoHomeHub</li>
                  </ul>
                </div>
            );
        }
        
        var zones
        if (this.state.zones.length > 0) {
            var self = this
            zones = this.state.zones.map(function(zone) {
                return <ZoneInfo errors={self.state.errors} ref={"zone" + zone.address} devices={self.state.devices} name={zone.name} description={zone.description} type="light" output="rgb" controller="FluxWIFI" address={zone.address} deviceId={zone.deviceId} key={zone.address} />
            })
        }

        var importBody
        if (!this.state.loading && this.state.devices.length > 0) {
            importBody = (
                <div>
                  <button className={"btn btn-primary" + (this.state.discovering ? " disabled" : "")}
                  onClick={this.discover}>Discover Zones</button>
                  <i className={"fa fa-spinner fa-spin discover" + (this.state.discovering ? "" : " hidden")}></i>
                  <h3 className={this.state.zones.length > 0 ? "" : " hidden"}>Zones</h3>
                  <p className={this.state.zones.length > 0 ? "" : " hidden"}>Click "Import" next to each zone you want to import</p>
                  {zones}
                </div>
            );
        }
        return (
            <div className="cmp-ImportFluxWIFI">
              {loading}
              {noDeviceBody}
              {importBody}
            </div>
        );
    }
});
module.exports = ImportFluxWIFI;