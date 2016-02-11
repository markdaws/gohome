var React = require('react');
var Zone = require('./Zone.jsx');
var Api = require('../utils/API.js');
var ZoneStore = require('../stores/ZoneStore.js');

var ZoneList = React.createClass({
    getInitialState: function() {
        return {
            zones: ZoneStore.getAll()
        }
    },
    
    componentDidMount: function() {
        ZoneStore.addChangeListener(this._onChange);
        Api.zoneLoadAll();
    },

    componentWillUnmount: function() {
        ZoneStore.removeChangeListener(this._onChange);
    },

    _onChange: function() {
        this.setState({ zones: ZoneStore.getAll() });
    },
    
    render: function() {
        //TODO: Add loading UI
        var zoneNodes = []
        this.state.zones.forEach(function(zone) {
            zoneNodes.push(
                <Zone id={zone.id} name={zone.name} type={zone.type} output={zone.output} key={zone.id}/>
            );
        })
        return (
            <div className="cmp-ZoneList row">
              {zoneNodes}
            </div>
        );
    }
});
module.exports = ZoneList;