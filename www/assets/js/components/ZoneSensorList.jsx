var ClassNames = require('classnames');
var React = require('react');
var ReactRedux = require('react-redux');
var ZoneControl = require('./ZoneControl.jsx');
var SensorMonitor = require('./SensorMonitor.jsx');
var ZoneActions = require('../actions/ZoneActions.js');
var Grid = require('./Grid.jsx');
var SensorMonitor = require('./SensorMonitor.jsx');
var ZoneSensorListGridCell = require('./ZoneSensorListGridCell.jsx');

var ZoneSensorList = React.createClass({
    render: function() {
        var lightZones = [];
        var shadeZones = [];
        var switchZones = [];
        var otherZones = [];
        var sensors = [];
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
        });

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
