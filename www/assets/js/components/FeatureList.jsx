var React = require('react');
var ReactRedux = require('react-redux');
var Grid = require('./Grid.jsx');
var Feature = require('../feature.js');
var FeatureCell = require('./FeatureCell.jsx');
var FeatureControl = require('./FeatureControl.jsx');
var FeatureInfo = require('./FeatureInfo.jsx');
var SystemActions = require('../actions/SystemActions.js');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'FeatureList',
    prefix: 'b-'
});
require('../../css/components/FeatureList.less')

//TODO: Need to get the correct state when we come out of edit mode, currently lost

var FeatureList = React.createClass({
    getInitialState: function() {
        this._monitorDataChanged = false;
        this._monitorId = '';
        this._connection = null;
        this._lastSubscribeId = 1;
        this._keepRefreshingConnection = true;
        this._retryDuration = 5000;
        this._monitorTimeout = 600;
        this._refreshTimeoutId = -1;
        this._expandedItems = [];
        this._monitorData = {
            features: {},
        };
        this.mergeFeatures(this.props.devices);
        return { editMode: false };
    },

    edit: function() {
        this.setState({ editMode: true });
    },

    endEdit: function() {
        this.setState({ editMode: false });
    },

    getDefaultProps: function() {
        return {
            zones: [],
            devices: [],
        };
    },

    mergeFeatures: function(devices) {
        devices.forEach(function(device) {
            (device.features || []).forEach(function(feature) {
                // If we have got monitor data back for an attribute we leave that otherwise
                // we set the attribute from the feature, that is like the blank initial state
                // that can be used to render the control before monitor data comes back. We
                // don't want to modify the original feature attribute since it is like a template
                // that gets copies, so we clone the attrs first before we use them

                var attrs = Feature.cloneAttrs(feature.attrs)
                var monitorFeature = this._monitorData.features[feature.id];
                if (!monitorFeature) {
                    this._monitorData.features[feature.id] = attrs;
                } else {
                    Object.keys(attrs).forEach(function(localId) {
                        if (!monitorFeature[localId]) {
                            monitorFeature[localId] = attrs[localId];
                        }
                    });
                }
            }.bind(this));
        }.bind(this));
    },

    componentWillReceiveProps: function(nextProps) {
        var shouldRefreshMonitor = false;
        if (nextProps.devices && (this.props.devices !== nextProps.devices)) {
            shouldRefreshMonitor = true;
            this.mergeFeatures(nextProps.devices);
        }
        if (shouldRefreshMonitor) {
            this._monitorDataChanged = true;
        }
    },

    shouldComponentUpdate: function(nextProps, nextState) {
        //TODO:
        return true;
        if (nextProps && (this.props.zones !== nextProps.zones)) {
            return true;
        }
        if (nextProps && (this.props.devices !== nextProps.devices)) {
            return true;
        }
        if (nextState && (nextState.editMode !== this.state.editMode)) {
            return true;
        }
        return false;
    },

    componentDidMount: function() {
        this.refreshMonitoring(this.props.devices);
    },

    componentWillUnmount: function() {
        if (this._monitorId !== '') {
            Api.monitorUnsubscribe(this._monitorId);
        }
        if (this._connection) {
            this._connection.close();
        }
        this._keepRefreshingConnection = false;
    },

    componentDidUpdate: function() {
        if (this._monitorDataChanged) {
            this.refreshMonitoring(this.props.devices);
        }
    },

    refreshMonitoring: function(devices) {
        this._monitorDataChanged = false;

        if (devices.length === 0) {
            return;
        }

        // Unsub from an old monitor group
        if (this._monitorId) {
            Api.monitorUnsubscribe(this._monitorId);
            this._monitorId = '';
        }

        // Just incase we are getting many refresh requests, due to re-rendering,
        // importing many devices, have a small delay before we actually fire a request
        clearTimeout(this._refreshTimeout);
        this._refreshTimeout = setTimeout(function() {
            this.refreshMonitoringInternal(devices);
        }.bind(this), 1000);
    },

    // do not call directly, call refreshMonitoring
    refreshMonitoringInternal: function(devices) {
        var monitorGroup = {
            timeoutInSeconds: this._monitorTimeout,
            featureIds: []
        };

        devices.forEach(function(device) {
            (device.features || []).forEach(function(feature) {
                monitorGroup.featureIds.push(feature.id);
            });
        });

        var subscribeId = ++this._lastSubscribeId;
        Api.monitorSubscribe(monitorGroup, function(err, data) {
            if (subscribeId !== this._lastSubscribeId) {
                // This is an old callback we subscribed to before the most recent, unsub
                Api.monitorUnsubscribe(data.monitorId);
                return
            }

            if (err != null) {
                setTimeout(function() {
                    if (this._keepRefreshingConnection) {
                        this.refreshMonitoring(this.props.zones, this.props.sensors);
                    }
                }.bind(this), this._retryDuration);
                return;
            }

            this._monitorId = data.monitorId;
            this.refreshWebSocket(this._monitorId);
        }.bind(this));
    },

    refreshWebSocket: function(monitorId) {
        if (this._connection) {
            this._connection.close();
            this._connection = null;
        }

        var conn = Api.monitorGroups(monitorId);
        conn.onopen = (function(evt) {
            function renew() {
                this._refreshTimeoutId = setTimeout(function() {
                    if (this._monitorId === '') {
                        return;
                    }
                    Api.monitorSubscribeRenew(this._monitorId);
                    renew.bind(this)();
                }.bind(this), parseInt(this._monitorTimeout * 1000 * 0.75, 10));
            }
            renew.bind(this)();
        }).bind(this);
        conn.onclose = function(evt) {
            clearTimeout(this._refreshTimeoutId);
            conn = null;
            this._connection = null;
            this._monitorId = '';

            setTimeout(function() {
                if (this._keepRefreshingConnection) {
                    this.refreshMonitoring(this.props.devices);
                }
            }.bind(this), 5000);
        }.bind(this);
        conn.onmessage = (function(evt) {
            var resp = JSON.parse(evt.data);
            Object.keys(resp.features || {}).forEach(function(featureId) {
                // Each feature can have multiple attributes, might only be some that
                // have updated, so here we merge the new results with the old results
                var attrs = this._monitorData.features[featureId] || {};
                var newAttrs = resp.features[featureId];
                Object.keys(newAttrs).forEach(function(localId) {
                    attrs[localId] = newAttrs[localId];
                });
                this._monitorData.features[featureId] = attrs;

                var cmp = this.refs['cell_feature_' + featureId];
                if (!cmp) {
                    return;
                }
                cmp.setAttrs(this._monitorData.features[featureId]);
            }.bind(this));

            this.updateExpandedItems();
        }).bind(this);
        this._connection = conn;
    },

    expanderMounted: function(content) {
        this._expandedItems.push(content);
        this.updateExpandedItems();
    },

    updateExpandedItems: function() {
        this._expandedItems.forEach(function(cmp) {
            cmp.monitorData(this._monitorData);
        }.bind(this));
    },

    expanderUnmounted: function(content) {
        this._expandedItems = this._expandedItems.filter(function(cmp) {
            return cmp !== content;
        });
    },

    attrChanged: function(feature, attr) {
        var data = {};
        data[attr.localId] = attr;

        Api.deviceSetFeatureAttrs(
            feature.deviceId,
            feature.id,
            data,
            function(err, data) {
                if (err) {
                    console.error(err);
                    //TODO:
                    return
                }
            }
        );
    },

    render: function() {
        var body, btns;

        if (this.state.editMode) {
            btns = (
                <div>
                    <div {...classes('buttons', 'editing', 'clearfix')}>
                        <button className="btn btn-default btnDone pull-right" onClick={this.endEdit}>
                            <i className="fa fa-times"></i>
                        </button>
                    </div>
                </div>
            );

            //TODO: Cache this, only update if features updated
            var features = [];
            this.props.devices.forEach(function(device) {
                (device.features || []).forEach(function(feature) {
                    // TODO:
                    // Don't support buttons right now
                    if (feature.type === Feature.Type.Button) {
                        return;
                    }
                    features.push(feature);
                })
            });
            features.sort(function(a, b) {
                return a.name.localCompare(b.name);
            });

            var featureCmps = [];
            features.forEach(function(feature) {
                featureCmps.push(
                    <div {...classes('feature-info')} key={feature.id}>
                        <FeatureInfo
                            readOnlyFields="id, deviceId"
                            key={feature.id}
                            feature={feature}
                            showSaveBtn={true}
                            updatedFeature={this.props.updatedFeature} />
                    </div>
                );
            }.bind(this));

            body = (
                <div>
                    {featureCmps}
                </div>
            );

        } else {
            var lightZones = [];
            var windowTreatments = [];
            var other = [];
            var sensors = [];

            this.props.devices.forEach(function(device) {
                (device.features || []).forEach(function(feature) {

                    var cmpFeature = {
                        key: 'feature_' + feature.id,
                        cell: <FeatureCell
                                  key={feature.id}
                                  ref={"cell_feature_" + feature.id}
                                  feature={feature} />,
                        content: <FeatureControl
                                     key={feature.id}
                                     id={feature.id}
                                     onAttrChanged={this.attrChanged}
                                     didMount={this.expanderMounted}
                                     willUnmount={this.expanderUnmounted}
                                     feature={feature} />
                    };

                    // For UI purposes we will group some of the more
                    // common features together
                    switch(feature.type) {
                        case Feature.Type.LightZone:
                            lightZones.push(cmpFeature)
                            break;
                        case Feature.Type.WindowTreatment:
                            windowTreatments.push(cmpFeature);
                            break;
                        case Feature.Type.Sensor:
                            sensors.push(cmpFeature)
                            break;
                        case Feature.Type.Button:
                            break;
                            //TODO: re-enable
                            //right now can't do anything with a button in the UI so hiding
                        default:
                            other.push(cmpFeature)
                    }
                }.bind(this));
            }.bind(this));

            btns = (
                <div {...classes('buttons', '', 'clearfix')}>
                    <button className="btn btn-default btnEdit pull-right" onClick={this.edit}>
                        <i className="fa fa-cog" aria-hidden="true"></i>
                    </button>
                </div>
            );

            body = (
                <div>
                    <div {...classes('grid-section', lightZones.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Lights</h2>
                        <Grid cells={lightZones} />
                    </div>
                    <div {...classes('grid-section', windowTreatments.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Window Treatments</h2>
                        <Grid cells={windowTreatments} />
                    </div>
                    <div {...classes('grid-section', sensors.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Sensors</h2>
                        <Grid cells={sensors} />
                    </div>
                    <div {...classes('grid-section', other.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Other</h2>
                        <Grid cells={other} />
                    </div>
                </div>
            );
        }

        return (
            <div {...classes()}>
                {btns}
                {body}
            </div>
        );
    }
});

function mapDispatchToProps(dispatch) {
    return {
        updatedFeature: function(featureJson) {
            dispatch(SystemActions.updatedFeature(featureJson));
        }
    }
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(FeatureList);
