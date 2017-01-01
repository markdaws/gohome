var React = require('react');

var FeaturePicker = React.createClass({
    getDefaultProps: function() {
        return {
            features: []
        };
    },

    getInitialState: function() {
        return {
            value: this.props.defaultId
        };
    },

    selected: function(evt) {
        this.setState({ value: evt.target.value });
        this.props.changed && this.props.changed(evt.target.value);
    },

    render: function() {
        var options = [];
        this.props.features.forEach(function(feature) {
            options.push(<option key={feature.id} value={feature.id}>{feature.name}</option>);
        });

        var noFeatures;
        var picker;
        if (this.props.features.length === 0) {
            noFeatures = <div>No features found.</div>
        } else {
            picker = (
                <select
                    disabled={this.props.disabled}
                    className="form-control"
                    onChange={this.selected}
                    value={this.state.value} >
                    <option value="">Select a feature...</option>
                    {options}
                </select>
            );
        }
        return (
            <div className="b-FeaturePicker">
                {picker}
                {noFeatures}
            </div>
        );
    }
});
module.exports = FeaturePicker;
