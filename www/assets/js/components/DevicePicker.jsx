var React = require('react');

var DevicePicker = React.createClass({
    getDefaultProps: function() {
        return {
            devices: []
        };
    },
    
    getInitialState: function() {
        return {
            value: this.props.defaultId
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
            var id = device.id;
            options.push(<option key={id} value={id}>{device.name}</option>);
        });
        return (
            <div className="cmp-DevicePicker">
                <select
                    disabled={this.props.disabled}
                    className="form-control"
                    onChange={this.selected}
                    value={this.state.value} >
                    <option value="">Select a device...</option>
                    {options}
                </select>
            </div>
        );
    }
});
module.exports = DevicePicker;
