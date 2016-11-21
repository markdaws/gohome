var React = require('react');
var ReactDOM = require('react-dom');
var Api = require('../utils/API.js');
var Attribute = require('../attribute.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'HueAttr',
    prefix: 'b-'
});
require('../../css/components/HueAttr.less')

var HueAttr = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.attr.value
        };
    },

    initSlider: function() {
        var sliders = $(ReactDOM.findDOMNode(this)).find('.b-HueAttr__slider');
        sliders.ColorPickerSliders({
            color: "rgb(36, 170, 242)",
            flat: true,
            swatches: false,
            order: {
                hsl: 1,
            }
        });
    },

    componentDidMount: function() {
        this._slider = this.initSlider();
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.attr && nextProps.attr != this.props.attr) {
            var newLevel = nextProps.attr.value;
            if (newLevel == null) {
                return;
            }
            this.setState({ value: newLevel });
            //this._slider && this._slider.set(Math.round(newLevel));
        }
    },

    setAttrs: function(attrs) {
        this.setState({ attrs: attrs });
    },

    render: function() {
        var val = '-';
        if (this.state.value != null) {
            val = this.state.value + '%';
        }

        var readOnly = this.props.attr.perms == Attribute.Perms.ReadOnly;
        return (
            <div {...classes('')}>
                <div {...classes('slider', readOnly ? 'read-only' : '')}></div>
            </div>
        );
    }
});
module.exports = HueAttr;
