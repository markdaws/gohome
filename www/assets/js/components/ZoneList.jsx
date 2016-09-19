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
        var zoneNodes = [];
        var zones = this.props.zones;
        zones.items.forEach(function(zone) {
            zoneNodes.push(
                <Zone id={zone.id} name={zone.name} type={zone.type} output={zone.output} key={zone.id}/>
            );
        })

        var loading;
        if (zones.loading) {
            loading = (
                <div className="spinnerWrapper">
                    <p>Loading Zones ...</p>
                    <i className="fa fa-spinner fa-spin"></i>
                </div>
            );
        }

        var error;
        if (zones.loadingErr) {
            error = <div>There was an error loading your zones. Please refresh the page.</div>;
        }

        var classNames = ClassNames({
            'cmp-ZoneList': true,
            'row': !zones.loadingErr
        });
        return (
            <div className={classNames}>
                {error}
                {loading}
                {zoneNodes}
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
