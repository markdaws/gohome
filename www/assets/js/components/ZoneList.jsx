var ClassNames = require('classnames');
var React = require('react');
var ReactRedux = require('react-redux');
var Zone = require('./Zone.jsx');
var ZoneActions = require('../actions/ZoneActions.js');

var ZoneList = React.createClass({
    render: function() {
        var lightZones = [];
        var shadeZones = [];
        var outletZones = [];
        var otherZones = [];
        this.props.zones.forEach(function(zone) {
            var cmpZone = (
                <Zone
                    id={zone.id}
                    name={zone.name}
                    type={zone.type}
                    output={zone.output}
                    key={zone.id}/>
            );

            switch(zone.type) {
                case 'light':
                    lightZones.push(cmpZone);
                    break;
                case 'shade':
                    shadeZones.push(cmpZone);
                    break;
                case 'outlet':
                    outletZones.push(cmpZone);
                    break;
                default:
                    otherZones.push(cmpZone);
                    break;
            }
        })

        return (
            <div className="cmp-ZoneList row">
                <div className="clearfix">
                <h2 className={ClassNames({ 'hidden': lightZones.length === 0})}>Lights</h2>
                {lightZones}
                </div>
                <div className="clearfix">
                <h2 className={ClassNames({ 'hidden': shadeZones.length === 0})}>Shades</h2>
                {shadeZones}
                </div>
                <div className="clearfix">
                <h2 className={ClassNames({ 'hidden': outletZones.length === 0})}>Outlets</h2>
                {outletZones}
                </div>
                <div className="clearfix">
                <h2 className={ClassNames({ 'hidden': otherZones.length === 0})}>Other Zones</h2>
                {otherZones}
                </div>
            </div>
        );
    }
});
module.exports = ZoneList;
