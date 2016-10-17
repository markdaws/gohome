var React = require('react');

var ZoneOutputPicker = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.output || 'continuous'
        };
    },

    componentDidMount: function() {
        // If a value wasn't passed in, raise a changed notification so callers
        // can set their value accordingly since we default to unknown
        if (this.state.value === 'continuous') {
            this.props.changed && this.props.changed(this.state.value);
        }
    },
    
    selected: function(evt) {
        this.setOutput(evt.target.value);
    },

    setOutput: function(output) {
        this.setState({ value: output });
        this.props.changed && this.props.changed(output);
    },
    
    render: function() {
        return (
            <div className="cmp-ZoneOutputPicker">
                <select
                    className="form-control"
                    onChange={this.selected}
                    value={this.state.value}>
                    <option value="unknown">Unknown</option>
                    <option value="continuous">Continuous</option>
                    <option value="binary">Binary</option>
                    <option value="rgb">RGB</option>
                </select>
            </div>
        );
    }
});
module.exports = ZoneOutputPicker;
