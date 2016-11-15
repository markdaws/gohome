var ClassNames = require('classnames');
var React = require('react');
var ReactRedux = require('react-redux');
var ZoneControl = require('./ZoneControl.jsx');
var SensorMonitor = require('./SensorMonitor.jsx');
var ZoneActions = require('../actions/ZoneActions.js');
var SensorActions = require('../actions/SensorActions.js');
var Grid = require('./Grid.jsx');
var SensorMonitor = require('./SensorMonitor.jsx');
var ZoneInfo = require('./ZoneInfo.jsx');
var SensorInfo = require('./SensorInfo.jsx');
var ZoneSensorListGridCell = require('./ZoneSensorListGridCell.jsx');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'ZoneSensorList',
    prefix: 'b-'
});
require('../../css/components/ZoneSensorList.less')

//TODO: Need to get the correct state when we come out of edit mode, currently lost

var ZoneSensorList = React.createClass({
    getInitialState: function() {
        this._monitorDataChanged = false;
        this._monitorId = '';
        this._connection = null;
        this._lastSubscribeId = 1;
        this._keepRefreshingConnection = true;
        this._retryDuration = 5000;
        this._monitorTimeout = 600;
        this._refreshTimeoutId = -1;
        this._monitorData = {
            zones: {},
            sensors: {}
        };
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
            sensors: [],
            zones: []
        };
    },
    
    componentWillReceiveProps: function(nextProps) {
        var shouldRefreshMonitor = false;
        if (nextProps.zones && (this.props.zones !== nextProps.zones)) {
            shouldRefreshMonitor = true;
        }
        if (nextProps.sensors && (this.props.sensors !== nextProps.sensors)) {
            shouldRefreshMonitor = true;
        }
        if (shouldRefreshMonitor) {
            this._monitorDataChanged = true;
        }
    },
    
    shouldComponentUpdate: function(nextProps, nextState) {
        if (nextProps && (this.props.zones !== nextProps.zones)) {
            return true;
        }
        if (nextProps && (this.props.sensors !== nextProps.sensors)) {
            return true;
        }
        if (nextState && (nextState.editMode !== this.state.editMode)) {
            return true;
        }
        return false;
    },

    componentDidMount: function() {
        this.refreshMonitoring(this.props.zones, this.props.sensors);
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
            this.refreshMonitoring(this.props.zones, this.props.sensors);
        }
    },

    refreshMonitoring: function(zones, sensors) {
        this._monitorDataChanged = false;

        if (zones.length === 0 && sensors.length === 0) {
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
            this.refreshMonitoringInternal(zones, sensors);
        }.bind(this), 1000);
    },

    // do not call directly, call refreshMonitoring
    refreshMonitoringInternal: function(zones, sensors) {
        var monitorGroup = {
            timeoutInSeconds: this._monitorTimeout,
            sensorIds: [],
            zoneIds: []
        };

        zones.forEach(function(zone) {
            monitorGroup.zoneIds.push(zone.id);
        });
        sensors.forEach(function(sensor) {
            monitorGroup.sensorIds.push(sensor.id);
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

        var conn = new WebSocket("ws://" + window.location.hostname + ":5000/api/v1/monitor/groups/" + monitorId);
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
                    this.refreshMonitoring(this.props.zones, this.props.sensors);
                }
            }.bind(this), 5000);
        }.bind(this);
        conn.onmessage = (function(evt) {
            var resp = JSON.parse(evt.data);
            Object.keys(resp.zones || {}).forEach(function(zoneId) {
                // Need to update our local data, since we can get back updates at any time
                // for any zone, we have to merge all the values into our one source of truth
                this._monitorData.zones[zoneId] = resp.zones[zoneId];
                
                var cmp = this.refs['cell_zone_' + zoneId];
                if (!cmp) {
                    return;
                }
                cmp.setLevel(this._monitorData.zones[zoneId]);
            }.bind(this));
            Object.keys(resp.sensors || {}).forEach(function(sensorId) {
                this._monitorData.sensors[sensorId] = resp.sensors[sensorId];
                
                var cmp = this.refs['cell_sensor_' + sensorId];
                if (!cmp) {
                    return;
                }
                cmp.setAttr(this._monitorData.sensors[sensorId]);
            }.bind(this));

            this._gridContent && this._gridContent.monitorData(this._monitorData);
        }).bind(this);
        this._connection = conn;
    },

    expanderMounted: function(content) {
        this._gridContent = content;
        this._gridContent.monitorData(this._monitorData);
    },

    expanderUnmounted: function() {
        this._gridContent = null;
    },
    
    render: function() {
        var body, btns;
        
        if (this.state.editMode) {
            btns = (
                <div>
                    <div {...classes('buttons', 'editing', 'clearfix')}>
                        <button className="btn btn-success btnDone pull-right" onClick={this.endEdit}>Done</button>
                    </div>
                </div>
            );

            var zones = this.props.zones.map(function(zone) {
                return (
                    <div {...classes('zone-info')} key={zone.id}>
                        <ZoneInfo
                            name={zone.name}
                            description={zone.description}
                            address={zone.address}
                            id={zone.id}
                            key={zone.id}
                            showSaveBtn={true}
                            readOnlyFields="deviceId, id"
                            deviceId={zone.deviceId}
                            type={zone.type}
                            devices={this.props.devices}
                            output={zone.output}
                            updatedZone={this.props.updatedZone} />
                    </div>
                );
            }.bind(this));

            var sensors = this.props.sensors.map(function(sensor) {
                return (
                    <div {...classes('sensor-info')} key={sensor.id}>
                        <SensorInfo
                            readOnlyFields="deviceId"
                            key={sensor.id}
                            name={sensor.name}
                            description={sensor.description}
                            address={sensor.address}
                            id={sensor.id}
                            attr={sensor.attr}
                            showSaveBtn={true}
                            deviceId={sensor.deviceId}
                            devices={this.props.devices}
                            updatedSensor={this.props.updatedSensor} />
                    </div>
                );
            }.bind(this));

            body = (
                <div>
                    {zones}
                    {sensors}
                </div>
            );

        } else {
            var lightZones = [];
            var shadeZones = [];
            var switchZones = [];
            var otherZones = [];
            var sensors = [];

            this.props.zones.forEach(function(zone) {
                var cmpZone = {
                    key: 'zones_' + zone.id,
                    cell: <ZoneSensorListGridCell
                              key={zone.id}
                              ref={"cell_zone_" + zone.id}
                              zone={zone} />,
                    content: <ZoneControl
                                 key={zone.id}
                                 id={zone.id}
                                 didMount={this.expanderMounted}
                                 willUnmount={this.expanderUnmounted}
                                 name={zone.name}
                                 type={zone.type}
                                 output={zone.output}
                                 key={zone.id}/>
                };
                switch(zone.type) {
                    case 'light':
                        lightZones.push(cmpZone);
                        break;
                    case 'shade':
                        shadeZones.push(cmpZone);
                        break;
                    case 'switch':
                        switchZones.push(cmpZone);
                        break;
                    default:
                        otherZones.push(cmpZone);
                        break;
                }
            }.bind(this));

            this.props.sensors.forEach(function(sensor) {
                var cmpSensor = {
                    key: 'sensor_' + sensor.id,
                    cell: <ZoneSensorListGridCell
                              key={sensor.id}
                              ref={"cell_sensor_" + sensor.id}
                              sensor={sensor} />,
                    content: <SensorMonitor
                                 id={sensor.id}
                                 didMount={this.expanderMounted}
                                 willUnmount={this.expanderUnmounted}
                                 key={sensor.id}
                                 sensor={sensor} />
                };
                sensors.push(cmpSensor);
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
                        <Grid cells={lightZones} expanderWillMount={this.zoneExpanderWillMount}/>
                    </div>
                    <div {...classes('grid-section', shadeZones.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Shades</h2>
                        <Grid cells={shadeZones} expanderWillMount={this.zoneExpanderWillMount}/>
                    </div>
                    <div {...classes('grid-section', switchZones.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Switches</h2>
                        <Grid cells={switchZones} expanderWillMount={this.zoneExpanderWillMount}/>
                    </div>
                    <div {...classes('grid-section', otherZones.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Other Zones</h2>
                        <Grid cells={otherZones} expanderWillMount={this.zoneExpanderWillMount}/>
                    </div>
                    <div {...classes('grid-section', sensors.length === 0 ? 'hidden' : '')}>
                        <h2 {...classes('grid-header')}>Sensors</h2>
                        <Grid cells={sensors} expanderWillMount={this.sensorExpanderWillMount}/>
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
        updatedZone: function(data) {
            dispatch(ZoneActions.updated(data));
        },
        updatedSensor: function(data) {
            dispatch(SensorActions.updated(data))
        }
    }
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(ZoneSensorList);
