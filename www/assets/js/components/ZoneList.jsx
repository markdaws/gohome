var ClassNames = require('classnames');
var React = require('react');
var ReactRedux = require('react-redux');
var Zone = require('./Zone.jsx');
var ZoneActions = require('../actions/ZoneActions.js');

var ZoneList = React.createClass({
    componentDidMount: function() {
        this.props.loadAllZones();
    },

    render: function() {
        var lightZones = [];
        var shadeZones = [];
        var otherZones = [];
        this.props.zones.items.forEach(function(zone) {
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

        var loading;
        if (this.props.zones.loading) {
            loading = (
                <div className="spinnerWrapper">
                    <p>Loading Zones ...</p>
                    <i className="fa fa-spinner fa-spin"></i>
                </div>
            );
        }

        var error;
        if (this.props.zones.loadingErr) {
            error = <div>There was an error loading your zones. Please refresh the page.</div>;
        }

        var classNames = ClassNames({
            'cmp-ZoneList': true,
            'row': !this.props.zones.loadingErr
        });
        return (
            <div className={classNames}>
                {error}
                {loading}
                <h2 className={ClassNames({ 'hidden': lightZones.length === 0 || loading})}>Lights</h2>
                {lightZones}
                <h2 className={ClassNames({ 'hidden': shadeZones.length === 0 || loading})}>Shades</h2>
                {shadeZones}
                <h2 className={ClassNames({ 'hidden': otherZones.length === 0 || loading})}>Other Zones</h2>
                {otherZones}
            </div>
        );
    }
});

function mapStateToProps(state) {
    return {
        zones: state.zones
    };
}

function mapDispatchToProps(dispatch) {
    return {
        loadAllZones: function() {
            dispatch(ZoneActions.loadAll());
        }
    }
}

var ZoneListContainer = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(ZoneList);
module.exports = ZoneListContainer;
