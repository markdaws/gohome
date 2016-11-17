var React = require('react');
var ReactRedux = require('react-redux');
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
            discovered: false,
            devices: null,
            scenes: null,
            uiFields: uiFields,
            errors: null
        };
    },

    discover: function() {
        this.setState({
            discovering: true,
            discovered: false,
            devices: null,
            scenes: null,
            errors: null
        });
        
        Api.discovererScanDevices(
            this.props.discoverer.id,
            this.state.uiFields,
            function(err, data) {
                if (err != null) {
                    this.setState({
                        discovering: false,
                        discovered: true,
                        errors: err
                    });
                    return;
                }
                this.setState({
                    discovering: false,
                    discovered: true,
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
        var deviceCount = 0;
        if (this.state.devices) {
            deviceCount = this.state.devices.length;
        }

        var uiFields = this.props.discoverer.uiFields.map(function(uiField) {
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

        var errors;
        if (this.state.errors) {
            errors = (
                <div {...classes('error')}>{this.state.errors.msg}</div>
            );
        }

        var importGroup;
        if (this.state.discovered && !this.state.errors) {
            importGroup = <ImportGroup
                    devices={this.state.devices}
                    createdDevice={this.props.importedDevice}
                    createdZones={this.props.importedZones}
                    createdSensors={this.props.importedSensors} />;
        }

        return (
            <div {...classes()}>
                <div {...classes('pre-import-instructions', this.props.discoverer.preScanInfo == '' ? 'hidden' : '')}>
                    {this.props.discoverer.preScanInfo}
                </div>
                <div>
                    {uiFields}
                </div>
                <div {...classes('discover')}>
                    <button {...classes('', '', (this.state.discovering ? 'disabled' : '') + ' btn btn-primary')}
                        onClick={this.discover}>Discover</button>
                    <i {...classes('spinner', this.state.discovering ? '' : 'hidden', 'fa fa-spinner fa-spin')}></i>
                </div>
                {errors}
                <div {...classes('import-group')}>
                    {importGroup}
                </div>
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
