var ClassNames = require('classnames');
var React = require('react');
var ReactRedux = require('react-redux');
var ZoneControl = require('./ZoneControl.jsx');
var SensorMonitor = require('./SensorMonitor.jsx');
var ZoneActions = require('../actions/ZoneActions.js');
var Grid = require('./Grid.jsx');
var SensorMonitor = require('./SensorMonitor.jsx');
var ZoneSensorListGridCell = require('./ZoneSensorListGridCell.jsx');
var Api = require('../utils/API.js');

var ZoneSensorList = React.createClass({
    getInitialState: function() {
        this._monitorDataChanged = false;
        this._monitorId = '';
        this._connection = null;
        this._lastSubscribeId = 1;
        this._keepRefreshingConnection = true;
        this._retryDuration = 5000;
        this._monitorTimeout = 10;
        this._refreshTimeoutId = -1;
        return {
            monitorData: null
        }
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
        if (this.props.zones !== nextProps.zones) {
            return true;
        }
        if (this.props.sensors !== nextProps.sensors) {
            return true;
        }
        if (this.state.monitorData != nextState.monitorData) {
            return true;
        }
        return false;
    },

    componentDidMount: function() {
        // Need this function since componentDidUpdate is not called on the initial render
        if (this._monitorDataChanged) {
            this.refreshMonitoring(this.props.zones, this.props.sensors);
        }
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
            if (err != null) {
                setTimeout(function() {
                    if (this._keepRefreshingConnection) {
                        this.refreshMonitoring(this.props.zones, this.props.sensors);
                    }
                }.bind(this), this._retryDuration);
                return;
            }

            if (subscribeId !== this._lastSubscribeId) {
                // This is an old callback we subscribed to before the most recent, unsub
                Api.monitorUnsubscribe(data.monitorId);
                return
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

            //TODO: Need to renew the subscription
            setTimeout(function() {
                if (this._keepRefreshingConnection) {
                    this.refreshMonitoring(this.props.zones, this.props.sensors);
                }
            }.bind(this), 5000);
        }.bind(this);
        conn.onmessage = (function(evt) {
            var resp = JSON.parse(evt.data);

            this.setState({ monitorData: resp });

            //TODO :Needed?
            /*
            Object.keys(resp.zones || {}).forEach(function(zoneId) {
                var cmp = this.refs['cell_zone_' + zoneId];
                if (cmp) {
                    cmp.setLevel(resp.zones[zoneId]);
                }
            }.bind(this));
            Object.keys(resp.sensors || {}).forEach(function(sensorId) {
                var cmp = this.refs['cell_sensor_' + sensorId];
                if (!cmp) {
                    return;
                }
                //TODO: Set attribute...cmp.setLevel(resp.sensors[sensorId].value);
            }.bind(this));
            */

            this._gridContent && this._gridContent.monitorData(resp);

        }).bind(this);
        this._connection = conn;
    },

    //TODO: This should come from the grid
    /*
    zoneExpanderMounted: function(content) {
        this._gridContent = content;
        this._gridContent.monitorData(this.state.monitorData);
    },*/
    
    render: function() {
        var lightZones = [];
        var shadeZones = [];
        var switchZones = [];
        var otherZones = [];
        var sensors = [];

        console.log('zone sensor list rendering');
        console.log(this.state.monitorData);
        
        this.props.zones.forEach(function(zone) {
            var level;
            if (this.state.monitorData) {
                level = this.state.monitorData.zones[zone.id];
                //console.log('zone: ' + zone.id + ', level: ' + level.value);
            }
            var cmpZone = {
                key: zone.id,
                cell: <ZoneSensorListGridCell
                          level={level}
                          key={zone.id}
                          ref={"cell_zone_" + zone.id}
                          zone={zone} />,
                content: <ZoneControl
                             id={zone.id}
                             level={level}
                //didMount={this.zoneExpanderMounted}
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
                cell: <ZoneSensorListGridCell ref={"cell_sensor_" + sensor.id} sensor={sensor} />,
                content: <SensorMonitor sensor={sensor} />
            };
            sensors.push(cmpSensor);
        });

        return (
            <div className="cmp-ZoneSensorList">
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': lightZones.length === 0})}>Lights</h2>
                    <Grid name="zone grid" cells={lightZones} expanderWillMount={this.zoneExpanderWillMount}/>
                </div>
                {/*
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': shadeZones.length === 0})}>Shades</h2>
                    <Grid cells={shadeZones} expanderWillMount={this.zoneExpanderWillMount}/>
                </div>
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': switchZones.length === 0})}>Switches</h2>
                    <Grid cells={switchZones} expanderWillMount={this.zoneExpanderWillMount}/>
                </div>
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': otherZones.length === 0})}>Other Zones</h2>
                    <Grid cells={otherZones} expanderWillMount={this.zoneExpanderWillMount}/>
                </div>
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': sensors.length === 0})}>Sensors</h2>
                    <Grid cells={sensors} />
                </div>*/}
            </div>
        );
    }
});
module.exports = ZoneSensorList;
