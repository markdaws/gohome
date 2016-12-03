var React = require('react');
var Feature = require('../feature.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SceneActionPicker',
    prefix: 'b-'
});

var SceneActionPicker = React.createClass({
    getDefaultProps: function() {
        return {
            excluded: {},
        }
    },

    getInitialState: function() {
        return {
            value: this.props.type || 'unknown'
        };
    },

    selected: function(evt) {
        this.setType(evt.target.value);
    },

    setType: function(type) {
        if (type === 'unknown') {
            return;
        }

        // Set back to unknown since we render a new command when this is selected
        this.setState({ value: 'unknown' });
        this.props.changed && this.props.changed(type);
    },

    render: function() {
        var types = [
            { str: "Add new action...", val:'unknown'}
        ];

        var excluded = this.props.excluded;
        if (!excluded[Feature.Type.Button]) {
            types.push({ str: "Button", val:Feature.Type.Button });
        }
        if (!excluded[Feature.Type.CoolZone]) {
            types.push({ str: "Cool Zone", val:Feature.Type.CoolZone });
        }
        if (!excluded[Feature.Type.HeatZone]) {
            types.push({ str: "Heat Zone", val:Feature.Type.HeatZone });
        }
        if (!excluded[Feature.Type.LightZone]) {
            types.push({ str: "Light Zone", val:Feature.Type.LightZone });
        }
        if (!excluded[Feature.Type.Sensor]) {
            types.push({ str: "Sensor", val:Feature.Type.Sensor });
        }
        if (!excluded[Feature.Type.Switch]) {
            types.push({ str: "Switch", val:Feature.Type.Switch });
        }
        if (!excluded[Feature.Type.Outlet]) {
            types.push({ str: "Outlet", val:Feature.Type.Outlet });
        }
        if (!excluded[Feature.Type.WindowTreatment]) {
            types.push({ str: "Window Treatment", val:Feature.Type.WindowTreatment });
        }

        var nodes = types.map(function(type) {
            return <option value={type.val} key={type.val}>{type.str}</option>;
        });
        return (
            <div {...classes()}>
                <select
                    className="form-control"
                    onChange={this.selected}
                    value={this.state.value}>
                {nodes}
              </select>
            </div>
        );
    }
});
module.exports = SceneActionPicker;
