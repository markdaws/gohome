var ClassNames = require('classnames');
var React = require('react');
var ReactDOM = require('react-dom');
var CssMixin = require('./CssMixin.jsx');
var Api = require('../utils/API.js');
var ClassNames = require('classnames');

var SensorMonitor = React.createClass({
    mixins: [CssMixin],
    getInitialState: function() {
        return {
            attr: null,
        }
    },

    componentDidMount: function() {
        this.props.didMount && this.props.didMount(this);
    },

    componentWillUnmount: function() {
        this.props.willUnmount && this.props.willUnmount();
    },

    monitorData: function(data) {
        if (!data || !data.sensors) {
            return;
        }
        var attr = data.sensors[this.props.id];
        if (attr == undefined) {
            return;
        }
        this.setState({ attr: attr });
    },

    render: function() {
        var val = '';
        if (this.state.attr) {
            val = this.state.attr.value;

            // If there is a states map, which gives value -> ui string then
            // use that string instead of the raw value
            var uiVal = this.state.attr.states && this.state.attr.states[val];
            if (uiVal) {
                val = uiVal;
            }
        }

        return (
            <div className="cmp-SensorMonitor">
                <div className="clearfix">
                    <div className="name pull-left">
                        {this.props.sensor.name}
                    </div>
                    <span className="value">
                        {val}
                    </span>
                </div>
            </div>
        );
    }
});
module.exports = SensorMonitor;
