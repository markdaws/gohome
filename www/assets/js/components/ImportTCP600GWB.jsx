var React = require('react');
var DeviceInfo = require('./DeviceInfo.jsx');

var ImportTCP600GWB = React.createClass({
    getInitialState: function() {
        return {
            location: "",
            locationFailed: false,
            discoveryInProgress: false,
            tokenInProgress: false,
            token: '',
            tokenError: false,
            tokenMissingAddress: false
        };
    },
    
    autoDiscover: function() {
        var self = this;
        this.setState({ discoveryInProgress: true });
        
        $.ajax({
            url: '/api/v1/discovery/TCP600GWB',
            dataType: 'json',
            cache: false,
            success: function(data) {
                self.setState({
                    location: data.location,
                    discoveryInProgress: false
                });
            },
            error: function(xhr, status, err) {
                self.setState({
                    locationFailed: true,
                    discoveryInProgress: false
                });
            }
        });
    },

    getToken: function() {
        var device = this.refs.devInfo.toJson();
        this.setState({
            tokenMissingAddress: false,
            tokenInProgress: true
        });
        
        if (device.address === '') {
            this.setState({
                tokenMissingAddress: true,
                tokenInProgress: false
            });
            return;
        }
        
        var self = this;
        $.ajax({
            url: '/api/v1/discovery/TCP600GWB/token?address=' + device.address,
            dataType: 'json',
            cache: false,
            success: function(data) {
                self.setState({
                    tokenInProgress: false,
                    token: data.token,
                    tokenError: data.unauthorized
                });
            },
            error: function(xhr, status, err) {
                self.setState({
                    tokenError: true,
                    tokenInProgress: false
                });
            }
        });
    },
    
    render: function() {
        return (
            <div className="cmp-ImportTCP600GWB">
              <p>Click to automatically retrieve the network address for this device</p>
              <div className="form-group has-error">
                <button className={"btn btn-primary" + (this.state.discoveryInProgress ? " disabled" : "")} onClick={this.autoDiscover}>Discover Address</button>
                <i className={"fa fa-spinner fa-spin" + (this.state.discoveryInProgress ? "" : " hidden")}></i>
                <span className={"help-block" + (this.state.locationFailed ? "" : " hidden")}>Error - Auto discovery failed, verify your TCP device is connected to the same network. If this continues to fail, use the official TCP app to get the device address</span>
              </div>
              <p>Click to retrive the security token. Only click this after pressing the "sync" button on your physical ConnectedByTCP hub</p>
              <div className="form-group has-error">
                <button className={"btn btn-primary" + (this.state.tokenInProgress ? " disabled" : "")} onClick={this.getToken}>Get Token</button>
                <i className={"fa fa-spinner fa-spin" + (this.state.tokenInProgress ? "" : " hidden")}></i>
                <span className={"help-block" + (this.state.tokenError ? "" : " hidden")}>Error - unable to get the token, make sure you press the physical "sync" button on the TCP hub device before clicking the "Get Token" button otherwise this will fail</span>
                <span className={"help-block" + (this.state.tokenMissingAddress ? "" : " hidden")}>Error - you must put a valid network address in the "Address" field first before clicking this button</span>
              </div>
              <DeviceInfo modelNumber="TCP600GWB" readOnlyFields="modelNumber" showToken="true" token={this.state.token} tokenError={this.state.tokenError} address={this.state.location} ref="devInfo"/>
            </div>
        )
    }
});
module.exports = ImportTCP600GWB;