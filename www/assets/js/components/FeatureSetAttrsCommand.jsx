var React = require('react');
var Feature = require('../feature.js');
var BEMHelper = require('react-bem-helper');
var DevicePicker = require('./DevicePicker.jsx');
var FeaturePicker = require('./FeaturePicker.jsx');
var FeatureControl = require('./FeatureControl.jsx');

var classes = new BEMHelper({
    name: 'FeatureSetAttrsCommand',
    prefix: 'b-'
});
require('../../css/components/FeatureSetAttrsCommand.less')

var FeatureSetAttrsCommand = React.createClass({
    getInitialState: function() {
        var deviceId = null;
        var features = [];
        var feature = null;

        var selectedFeatureId = this.props.command.attributes && this.props.command.attributes.id;
        if (selectedFeatureId) {
            // We have a feature ID, so we need to load the original values
            // the user saved, vs this being a new command the user hasn't
            // chosen any values for yet

            // Find the device the feature belongs to
            // TODO: Augment on server side
            for (var i=0; i<this.props.devices.length; ++i) {
                var device = this.props.devices[i];
                for (var j=0; j<device.features.length; ++j) {
                    if (device.features[j].id === selectedFeatureId) {
                        deviceId = device.id;
                        feature = device.features[j];
                        features = device.features;
                        break;
                    }
                }
                if (feature) {
                    break;
                }
            }
        } else {
            // If we only have one device, select it by default
            if (this.props.devices.length === 1) {
                var device = this.props.devices[0];
                deviceId = device.id;

                // If we only have one feature that matches, we can just
                // select that automatically
                features = this.filterFeatures(device.id);
                if (features.length === 1) {
                    feature = features[0];
                }
            }
        }

        return {
            deviceId: deviceId,
            features: features,
            feature: feature
        };
    },

    filterFeatures: function(deviceId) {
        var features = [];
        var featureType = this.props.command.attributes.type;
        for (var i=0; i<this.props.devices.length; ++i) {
            var device = this.props.devices[i];

            if (device.id === deviceId) {
                (device.features || []).forEach(function(feature) {
                    if (feature.type === featureType) {
                        features.push(feature);
                    }
                });
                break;
            }
        }
        return features;
    },

    devicePickerChanged: function(deviceId) {
        var features = this.filterFeatures(deviceId);
        this.setState({
            deviceId: deviceId,
            features: features
        });
    },

    featurePickerChanged: function(featureId) {
        var feature;
        var features = this.state.features;
        for (var i=0; i<features.length; ++i) {
            if(features[i].id === featureId) {
                feature = features[i];
                break;
            }
        }

        this.setState({ feature: feature });
    },

    attrChanged: function(feature, attr) {
        this.props.onAttrChanged && this.props.onAttrChanged(feature, attr);
    },

    getSettings: function() {
        var featureControl = this.refs.featureControl;

        return {
            feature: this.state.feature,
            modifiedAttrs: featureControl.modifiedAttrs()
        };
    },

    render: function() {
        var command = this.props.command;
        var body, header;

        var featurePicker;
        var noFeatures;
        if (this.state.features.length > 0) {
            featurePicker = (
                <div {...classes('feature-picker')}>
                    <label>Feature:</label>
                    <FeaturePicker
                        disabled={this.props.command.id}
                        features={this.state.features}
                        defaultId={(this.state.feature || {}).id}
                        changed={this.featurePickerChanged}/>
                </div>
            );
        } else if (this.state.deviceId) {
            // User has chosen a device, but it doesn't have any features that
            // match the feature type
            noFeatures = <div>No matching feature found on this device</div>
        }

        var featureControl;
        if (this.state.feature) {
            featureControl = (
                <FeatureControl
                    hideReadOnlyAttrs={true}
                    ref="featureControl"
                    key={this.state.feature.id}
                    feature={this.state.feature}
                    onAttrChanged={this.attrChanged}
                    attrs={Feature.cloneAttrs(command.attributes.attrs)}
                />
            );
        }

        body = (
            <div>
                <div {...classes('device-picker')}>
                    <label>Device:</label>
                    <DevicePicker
                        disabled={this.props.command.id}
                        devices={this.props.devices}
                        defaultId={this.state.deviceId}
                        changed={this.devicePickerChanged}/>
                </div>
                {featurePicker}
                {noFeatures}
                {featureControl}
            </div>
        );

        return (
            <div {...classes()}>
                {body}
            </div>
        );
    }
});
module.exports = FeatureSetAttrsCommand;
