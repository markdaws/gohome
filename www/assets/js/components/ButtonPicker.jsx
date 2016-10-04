var React = require('react');

var ButtonPicker = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.buttonId || ''
        };
    },

    selected: function(evt) {
        this.setState({ value: evt.target.value });
        this.props.changed && this.props.changed(evt.target.value);
    },
    
    render: function() {
        var options = [];
        this.props.buttons.forEach(function(button) {
            options.push(<option key={button.id} value={button.id}>{button.fullName}</option>);
        });
        return (
            <div className="cmp-ButtonPicker">
                <select
                    disabled={this.props.disabled}
                    className="form-control"
                    defaultValue={this.props.buttonId}
                    onChange={this.selected}
                    value={this.state.value}>
                    <option value="">Select a Button...</option>
                    {options}
                </select>
            </div>
        );
    }
});
module.exports = ButtonPicker;
