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
    render: function() {
        var lightZones = [];
        var shadeZones = [];
        var switchZones = [];
        var otherZones = [];
        var sensors = [];

        //TODO: In wrong place - only do once, mounted added and re-rendered multiple times...
        var monitorGroup = {
            timeoutInSeconds: 200,
            sensorIds: [],
            zoneIds: []
        };

        this.props.zones.forEach(function(zone) {

            var cmpZone = {
                cell: <ZoneSensorListGridCell zone={zone} />,
                content: <ZoneControl
                             id={zone.id}
                             name={zone.name}
                             type={zone.type}
                             output={zone.output}
                             key={zone.id}/>
            };
                
            monitorGroup.zoneIds.push(zone.id);
            
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
        });

        this.props.sensors.forEach(function(sensor) {
            var cmpSensor = {
                cell: <ZoneSensorListGridCell sensor={sensor} />,
                content: <SensorMonitor sensor={sensor} />
            };
            sensors.push(cmpSensor);

            monitorGroup.sensorIds.push(sensor.id)
        });

        //TODO: Remove
        Api.monitorSubscribe(monitorGroup, function(err, data) {
            if (err != null) {
                console.log('failed to sub to monitor');
                console.log(err);
                return;
            }
            console.log('subscribed to monitor');
            console.log(data);

            reconnect(data.monitorId);
        });

        function reconnect(monitorId) {
            /*var oldConn = this.state.conn;
            if (oldConn) {
                oldConn.close();
            }*/

            var conn = new WebSocket("ws://" + window.location.hostname + ":5000/api/v1/monitor/groups/" + monitorId);
            var self = this;
            conn.onopen = function(evt) {
                /*
                self.setState({
                    connectionStatus: 'connected'
                });*/
            };
            conn.onclose = function(evt) {
                conn = null;
                /*
                self.setState({
                    conn: null,
                    items: [],
                    connectionStatus: 'disconnected'
                });*/
            };
            conn.onmessage = function(evt) {
                var item = JSON.parse(evt.data);
                console.log('got monitor message');
                console.log(item);
            };
            /*
            this.setState({
                conn: conn,
                connectionStatus: 'connecting'
            });*/
        }

        return (
            <div className="cmp-ZoneSensorList">
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': lightZones.length === 0})}>Lights</h2>
                    <Grid cells={lightZones} />
                </div>
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': shadeZones.length === 0})}>Shades</h2>
                    <Grid cells={shadeZones} />
                </div>
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': switchZones.length === 0})}>Switches</h2>
                    <Grid cells={switchZones} />
                </div>
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': otherZones.length === 0})}>Other Zones</h2>
                    <Grid cells={otherZones} />
                </div>
                <div className="clearfix">
                    <h2 className={ClassNames({ 'hidden': sensors.length === 0})}>Sensors</h2>
                    <Grid cells={sensors} />
                </div>
            </div>
        );
    }
});
module.exports = ZoneSensorList;
