var React = require('react');
var Zone = require('./Zone.jsx');

var ZoneList = React.createClass({
    render: function() {
        //TODO: Add loading
        var self = this;
        var zoneNodes = Object.keys(this.props.zones).map(function(id) {
            var zone = self.props.zones[id];
            return (
                <Zone id={zone.id} name={zone.name} type={zone.type} output={zone.output} key={id}/>
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