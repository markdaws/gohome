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

    render: function() {
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
