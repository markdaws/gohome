var React = require('react');
var ReactDOM = require('react-dom');
var Api = require('../utils/API.js');
var Attribute = require('../attribute.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'TempAttr',
    prefix: 'b-'
});
require('../../css/components/TempAttr.less')

var TempAttr = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.attr.value
        };
    },

    initSlider: function(step, min, max) {
        var sliders = $(ReactDOM.findDOMNode(this)).find('.b-TempAttr__slider');
        if (!sliders || sliders.length === 0) {
            return null;
        }

        var slider = sliders[0];
        noUiSlider.create(
            slider,
            {
                connect: [true, false],
                start: 0,
                animate: false,
                step: step,
                orientation: 'horizontal',
                range: {
                    min: min,
                    max: max
                }
            });
        slider.noUiSlider.set(this.state.value);
        slider.noUiSlider.on('slide', this.sliderChanged.bind(this, slider.noUiSlider));
        slider.noUiSlider.on('change', this.sliderEnd.bind(this, slider.noUiSlider));

        return slider.noUiSlider;
    },

    sliderChanged: function(slider) {
        this.setState({ value: parseInt(slider.get(), 10) });
    },

    sliderEnd: function(slider) {
        this.props.onTempChanged && this.props.onTempChanged(this.props.attr, parseInt(slider.get(), 10));
    },

    componentDidMount: function() {
        this._slider = this.initSlider(
            this.props.attr.step,
            this.props.attr.min,
            this.props.attr.max
        );
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.attr && nextProps.attr != this.props.attr) {
            var newLevel = nextProps.attr.value;
            if (newLevel == null) {
                return;
            }
            this.setState({ value: newLevel });
            this._slider && this._slider.set(Math.round(newLevel));
        }
    },

    setAttrs: function(attrs) {
        this.setState({ attrs: attrs });
    },

    render: function() {
        var val = '-';
        if (this.state.value != null) {
            val = this.state.value + '°F';
        }

        var readOnly = this.props.attr.perms == Attribute.Perms.ReadOnly;
        return (
            <div {...classes('', '', 'clearfix')}>
                <div {...classes('name')}>{this.props.attr.name}</div>
                <div {...classes('slider', readOnly ? 'read-only' : '')}></div>
                <span {...classes('value')}>{val}</span>
            </div>
        );
    }
});
module.exports = TempAttr;
