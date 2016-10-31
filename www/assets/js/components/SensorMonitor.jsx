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
            value: -1,
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
        var val = data.sensors[this.props.id];
        if (val == undefined) {
            return;
        }
        this.setState({ value: val });
    },

    render: function() {
        console.log(this.state.value);
        return (
            <div className="cmp-SensorMonitor">
                <div className="clearfix">
                    <div className="name pull-left">
                        {this.props.sensor.name}
                    </div>
                </div>
            </div>
        );
    }
});
module.exports = SensorMonitor;
