var React = require('react');
var ReactDOM = require('react-dom');
var Api = require('../utils/API.js');
var Attribute = require('../attribute.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'HSLAttr',
    prefix: 'b-'
});
require('../../css/components/HSLAttr.less')

var HSLAttr = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.attr.value
        };
    },

    initSlider: function() {
        var self = this;
        var sliders = $(ReactDOM.findDOMNode(this)).find('.b-HSLAttr__slider');

        sliders.ColorPickerSliders({
            color: 'hsl(0, 100%, 50%)',
            flat: true,
            swatches: false,
            order: {
                hsl: 1,
            },
            onchange: function(container, color) {
                var hsl = color.tiny.toHsl();
                if (self._ignoreChange || (hsl.h === 0 && hsl.s === 1 && hsl.l === 0.5)) {
                    // If we got an update from the API, shouldn't fire a change event here
                    // since we only want to do that when the user changes the slider. No way
                    // to distinguish with this control
                    self._ignoreChange = false;

                    // This is the default value, ignore since this fires on load an we can't
                    // distinguish between it and the user moving the slider. There is a bug
                    // here that the user can't set 0,0,0 since we ignore it now.
                    return;
                }

                // this fires many times as they are sliding, let them stop moving before we
                // send the commands
                clearTimeout(self._timeoutId);
                self._timeoutId = setTimeout(function() {
                    var hsl = color.tiny.toHslString();
                    self.setState({ value: hsl });
                    self.props.onChanged && self.props.onChanged(self.props.attr, hsl);
                }, 500);
            }
        });
        return sliders;
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

            // This should not trigger an update, just needed to update UI.
            this._ignoreChange = true
            this._slider && this._slider.trigger("colorpickersliders.updateColor", newLevel);
        }
    },

    setAttrs: function(attrs) {
        this.setState({ attrs: attrs });
    },

    render: function() {
        var readOnly = this.props.attr.perms == Attribute.Perms.ReadOnly;
        return (
            <div {...classes('')}>
                <div {...classes('slider', readOnly ? 'read-only' : '')}></div>
            </div>
        );
    }
});
module.exports = HSLAttr;
