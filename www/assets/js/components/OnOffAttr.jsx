var React = require('react');
var ReactDOM = require('react-dom');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'OnOffAttr',
    prefix: 'b-'
});
require('../../css/components/OnOffAttr.less')

var OnOffAttr = React.createClass({
    getInitialState: function() {
        console.log('got state');
        return {
            value: this.props.attr.value
        };
    },

    initSwitch: function(slider) {
        var sw = $($(ReactDOM.findDOMNode(this)).find('.switch-indeterminate')[0]);

        var options = {
            onText: 'On',
            offText: 'Off'
        };
        if (this.state.value != null) {
            options.state = this.state.value === 2;
        }

        sw.bootstrapSwitch(options);
        sw.on('switchChange.bootstrapSwitch', function(event, state) {
            this.toggleOn(slider);
        }.bind(this));

        // For some reason have to set this explicitly, not working in options
        sw.bootstrapSwitch('state', options.state, true);
        return sw;
    },

    toggleOn: function(slider) {
        var newValue;
        var off = 1;
        var on = 2;
        if (this.state.value == null) {
            newValue = off;
        } else if (this.state.value === on) {
            newValue = off
        } else {
            newValue = on;
        }
        this.setState({ value: newValue });
        this.props.onToggle && this.props.onToggle(this.props.attr, newValue);
    },

    componentDidMount: function() {
        this._switch = this.initSwitch(this._slider);
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.attr && nextProps.attr != this.props.attr) {
            var state = null;
            if (nextProps.attr.value != null) {
                state = nextProps.attr.value === 2;
            }
            this.setState({ value: nextProps.attr.value });
            this._switch && this._switch.bootstrapSwitch('state', state, true);
        }
    },

    render: function() {
        return (
            <div {...classes('')}>
                <input
                    className="switch-indeterminate"
                    type="checkbox"
                    defaultChecked={true}
                    data-indeterminate="true"></input>
            </div>
        );
    }
});
module.exports = OnOffAttr;
