var React = require('react');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');
var DeviceCell = require('./DeviceCell.jsx');
var FeatureCell = require('./FeatureCell.jsx');
var DeviceInfo = require('./DeviceInfo.jsx');
var Grid = require('./Grid.jsx');
var FeatureInfo = require('./FeatureInfo.jsx');
var SaveBtn = require('./SaveBtn.jsx');

var classes = new BEMHelper({
    name: 'ImportGroup',
    prefix: 'b-'
});
require('../../css/components/ImportGroup.less')

var ImportGroup = React.createClass({
    getInitialState: function() {
        this._unselected = {};
        this._savedDevices = {};
        this._savedFeatures = {};

        var features = [];
        (this.props.devices || []).forEach(function(device) {
            (device.features || []).forEach(function(feature) {
                features.push(feature);
            });
        });

        return {
            devices: this.props.devices,
            features: features,
            featureErrors: {},
            deviceErrors: {},
            saveButtonStatus: ''
        };
    },

    updateUnselected: function(id, checked) {
        if (checked) {
            delete this._unselected[id];
        } else {
            this._unselected[id] = true;
        }
    },

    deviceChkBxChanged: function(deviceId, checked) {
        this.updateUnselected(deviceId, checked);
    },

    featureChkBxChanged: function(featureId, checked) {
        this.updateUnselected(featureId, checked);
    },

    _deviceChanged: function(cmp) {
        var device = cmp.toJson();
        var devices = this.state.devices.map(function(dev) {
            if (dev.id === device.id) {
                return device;
            }
            return dev;
        });
        this.setState({
            devices: devices,
            saveButtonStatus: ''
        });
    },

    _featureChanged: function(cmp) {
        var feature = cmp.toJson();
        var features = this.state.features.map(function(f) {
            if (f.id === feature.id) {
                return feature;
            }
            return f;
        });
        this.setState({
            features: features,
            saveButtonStatus: ''
        });
    },

    importClicked: function() {
        this.refs["devicegrid"].closeExpander();
        this.refs["featuregrid"].closeExpander();
        this.setState({
            deviceErrors: {},
            featureErrors: {},
        });

        // When we are saving devices, we need to sort the devices so that
        // we save devices that don't have a hub before devices that do have
        // a hub, because we can't save a device that points to a hub that we
        // haven't potentially saved yet. Don't support hubs having hubs right
        // now otherwise we would need a dependency graph traversal instead of
        // a simple sort
        var sortedDevices = ([].concat(this.state.devices)).sort(function(x, y) {
            var xHasHub = x.hubId !== '';
            var yHasHub = y.hubId !== '';

            if ((xHasHub && yHasHub) || (!xHasHub && !yHasHub)) {
                return 0;
            } else if (xHasHub && !yHasHub) {
                return 1;
            } else {
                return -1;
            }
        })

        function saveDevice(devIndex) {
            if (devIndex >= sortedDevices.length) {
                this.setState({ saveButtonStatus: 'success' });
                return;
            }

            var device = sortedDevices[devIndex];

            //TODO: Object.assign

            // We grouped all of the features together, now we need to break them
            // up in to their respective devices so they can be saved
            device.features = this.state.features.filter(function(feature) {
                return feature.deviceId === device.id &&
                       !feature.isDupe &&
                       !this._unselected[feature.id] &&
                       !this._savedFeatures[feature.id];
            }.bind(this));

            if (device.isDupe || this._savedDevices[device.id]) {
                var saveFeature = (index) => {
                    if (index >= device.features.length) {
                        // Saved all the features, move on to the next device
                        saveDevice.bind(this)(devIndex + 1);
                        return;
                    }

                    var feature = device.features[index];
                    if (this._savedFeatures[feature.id]) {
                        // Already saved, can skip
                        saveFeature(index + 1);
                        return
                    }

                    Api.featureCreate(device.id, feature, function(err, featureData) {
                        if (err) {
                            var featureErrs = {};
                            featureErrs[feature.id] = err.validation.errors[feature.id];
                            this.setState({
                                saveButtonStatus: 'error',
                                featureErrors: featureErrs
                            });
                            return
                        }

                        this.props.createdFeature(feature);
                        this._savedFeatures[feature.id] = true;
                        saveFeature(index + 1);
                    }.bind(this));
                }
                saveFeature(0);
            } else {
                Api.deviceCreate(device, function(err, deviceData) {
                    if (err) {
                        var featureErrs = {};
                        (device.features || []).forEach(function(feature) {
                            if (err.validation.errors[feature.id]) {
                                featureErrs[feature.id] = err.validation.errors[feature.id];
                            }
                        });
                        this.setState({
                            saveButtonStatus: 'error',
                            featureErrors: featureErrs,
                            deviceErrors: err.validation.errors[device.id] || {}
                        });
                        return;
                    }

                    // Let callers know the device has been saved
                    this.props.createdDevice(deviceData);
                    //TODO: Feature

                    this._savedDevices[device.id] = true;
                    (device.features || []).forEach(function(feature) {
                        this._savedFeatures[feature.id] = true;
                    }.bind(this));

                    // Move on to the next device
                    saveDevice.bind(this)(devIndex + 1);
                }.bind(this));
            }
        }
        saveDevice.bind(this)(0);
    },

    getDeviceById: function(id) {
        var devices = this.state.devices;
        for (var i=0; i<devices.length; ++i) {
            if (devices[i].id === id) {
                return devices[i];
            }
        }
        return null;
    },

    render: function() {
        var devices = [];
        var features = [];

        (this.state.devices || []).forEach(function(device) {
            if (device.isDupe) {
                return;
            }

            var err = this.state.deviceErrors[device.id];
            var cell = {
                key: device.id,
                cell: <DeviceCell
                          key={"devicecell-" + device.id}
                          id={device.id}
                          showCheckbox={false}
                          chkBxChanged={this.deviceChkBxChanged}
                          hasError={err != null}
                          hasSuccess={this._savedDevices[device.id]}
                          device={device} />,
                content: <DeviceInfo
                             {...device}
                             key={"deviceinfo-" + device.id}
                             ref={"deviceinfo-" + device.id}
                             readOnlyFields="id"
                             errors={err}
                             changed={this._deviceChanged} />
            };
            devices.push(cell);
        }.bind(this));

        (this.state.features || []).forEach(function(feature) {
            if (feature.isDupe) {
                return;
            }

            var err = this.state.featureErrors[feature.id];
            var device = this.getDeviceById(feature.deviceId);
            var cmpFeature = {
                key: 'feature_' + feature.id,
                cell: <FeatureCell
                          id={feature.id}
                          showCheckbox={true}
                          showLevel={false}
                          chkBxChanged={this.featureChkBxChanged}
                          hasError={err != null}
                          hasSuccess={this._savedFeatures[feature.id]}
                          key={feature.id}
                          ref={"cell_feature_" + feature.id}
                          feature={feature} />,
                content: <FeatureInfo
                             feature={feature}
                             readOnlyFields="id, deviceId"
                             key={feature.id}
                             errors={err}
                             deviceId={device.id}
                             devices={[ device ]}
                             changed={this._featureChanged} />

            };
            features.push(cmpFeature);
        }.bind(this));

        var hasNonDupes = devices.length > 0 || features.length > 0;
        var body;
        if (hasNonDupes) {
            body = (
                <div>
                    <div {...classes('import', '', 'pull-right clearfix')}>
                        <SaveBtn
                            clicked={this.importClicked}
                            text="Import" status={this.state.saveButtonStatus} />
                    </div>
                    <div {...classes('devices', devices.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('header')}>Devices</h2>
                        <Grid ref="devicegrid" key="devicegrid" cells={devices} />
                    </div>
                    <div {...classes('features', features.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('header')}>Features</h2>
                        <Grid ref="featuregrid" key="featuregrid" cells={features} />
                    </div>
                    <div style={{ clear: 'both' }}></div>
                </div>
            );
        } else {
            body = (
                <div {...classes('no-new')}>
                    No new devices or features were found. Items previously imported will not be shown again
                    unless you delete them from the system.
                </div>
            );
        }
        return (
            <div {...classes('','', 'clearfix')}>
                {body}
            </div>
        );
    }
});
module.exports = ImportGroup;
