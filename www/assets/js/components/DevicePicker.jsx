var React = require('react');

module.exports = React.createClass({
    getInitialState: function() {
        return {
            value: ''
        };
    },

    //TODO: If only one item in the list, select by default on load
    //TODO: if output or type is unknown need to update zone control to be
    //able to handle those values
    selected: function(evt) {
        this.setState({ value: evt.target.value });
        this.props.changed && this.props.changed(evt.target.value);
    },
    
    render: function() {
        var options = [];
        this.props.devices.forEach(function(device) {
            options.push(<option key={device.id} value={device.id}>{device.name}</option>);
        });
        return (
            <div className="cmp-DevicePicker">
              <select className="form-control" onChange={this.selected} value={this.state.value}>
                <option value="">Select a device...</option>
                {options}
              </select>
            </div>
        );
    }
});
