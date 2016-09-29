var ClassNames = require('classnames');
var React = require('react');
var ReactRedux = require('react-redux');
var Zone = require('./Zone.jsx');
var ZoneActions = require('../actions/ZoneActions.js');

var ZoneList = React.createClass({
    render: function() {
        var lightZones = [];
        var shadeZones = [];
        var otherZones = [];
        this.props.zones.forEach(function(zone) {
            var cmpZone = <Zone id={zone.id} name={zone.name} type={zone.type} output={zone.output} key={zone.id}/>;

            switch(zone.type) {
                    //TODO: Put in enum somewhere
                case 'light':
                    lightZones.push(cmpZone);
                    break;
                case 'shade':
                    shadeZones.push(cmpZone);
                    break;
                default:
                    otherZones.push(cmpZone);
                    break;
            }
        })

        return (
            <div className="cmp-ZoneList row">
                <h2 className={ClassNames({ 'hidden': lightZones.length === 0})}>Lights</h2>
                {lightZones}
                <h2 className={ClassNames({ 'hidden': shadeZones.length === 0})}>Shades</h2>
                {shadeZones}
                <h2 className={ClassNames({ 'hidden': otherZones.length === 0})}>Other Zones</h2>
                {otherZones}
            </div>
        );
    }
});
module.exports = ZoneList;
