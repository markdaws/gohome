var ClassNames = require('classnames');
var React = require('react');
var ReactDOM = require('react-dom');
var CssMixin = require('./CssMixin.jsx');
var Api = require('../utils/API.js');
var ClassNames = require('classnames');

var ZoneControl = React.createClass({
    mixins: [CssMixin],
    getInitialState: function() {
        return {
            value: -1,
            r: 0,
            g: 0,
            b: 0
        }
    },

    initSlider: function() {
        var sliders = $(ReactDOM.findDOMNode(this)).find('.slider');
        if (!sliders || sliders.length === 0) {
            return null;
        }

        var slider = sliders[0];
        noUiSlider.create(
            slider,
            {
                connect: [true, false],
                start: 0,
                step: 1,
                orientation: 'horizontal',
                range: {
                    min: 0,
                    max: 100
                }
            });
        slider.noUiSlider.on('slide', this.sliderChanged.bind(this, slider.noUiSlider));
        slider.noUiSlider.on('change', this.sliderEnd.bind(this, slider.noUiSlider));

        return slider.noUiSlider;
    },

    initSwitch: function(slider) {
        var sw = $($(ReactDOM.findDOMNode(this)).find('.switch-indeterminate')[0]);
        sw.bootstrapSwitch({
            onText: this.props.type === 'shade' ? 'open' : 'on',
            offText: this.props.type === 'shade' ? 'close' : 'off',
        });

        sw.on('switchChange.bootstrapSwitch', function(event, state) {
            this.toggleOn(slider);
        }.bind(this));
    },

    initRGB: function() {
        var wrapper = $(ReactDOM.findDOMNode(this)).find('.rgbWrapper');
        if (!wrapper || wrapper.length === 0) {
            return;
        }

        wrapper.colorpicker({
            format: 'rgb',
            container: true,
            inline: true
        });

        var timeoutId;
        wrapper.colorpicker().on('changeColor', function(evt) {
            // Limit the number of events we send as the user is moving around
            // on the RGB UI, only send when they pause
            clearTimeout(timeoutId);
            timeoutId = setTimeout(function() {
                var rgb = evt.color.toRGB();
                this.setValue(
                    'setLevel',
                    0,
                    parseInt(rgb.r),
                    parseInt(rgb.g),
                    parseInt(rgb.b),
                    function(err) {
                        if (err) {
                            //TODO: err
                            console.error(err);
                        }
                    }
                );
            }.bind(this), 100);
        }.bind(this));
    },

    sliderChanged: function(slider) {
        this.setState({ value: parseInt(slider.get(), 10) });
    },

    sliderEnd: function(slider) {
        this.setValue('setLevel', parseInt(slider.get(), 10), 0, 0, 0, function(err) {
            if (err) {
                //TODO:
                console.error(err);
            }
        });        
    },

    componentDidMount: function() {
        var level = this.props.getZoneLevel(this.props.id);
        if (level) {
            this.setState({
                value: level.value,
                r : level.r,
                g: level.g,
                b: level.b
            });
        }
        var slider = this.initSlider();
        this.initSwitch(slider);
        this.initRGB();
    },

    setValue: function(cmd, value, r, g, b, callback) {
        this.setState({
            value: value,
            r: r,
            g: g,
            b: b
        });
        this.send({
            cmd: cmd,
            value: parseFloat(value),
            r: r,
            g: g,
            b: b
        }, callback);
    },

    toggleOn: function(slider) {
        var cmd, level;
        if (this.state.value !== 0 || this.state.r !== 0 || this.state.g !== 0 || this.state.b !== 0) {
            cmd = 'turnOff';
            level = 0;
        } else {
            cmd = 'turnOn';
            level = 100;
        }
        
        slider && slider.set(level);
        this.setValue(cmd, level, 0, 0, 0, function(err) {
            if (err) {
                //TODO: error
                console.error(err);
            }
        });
    },

    send: function(data, callback) {
        Api.zoneSetLevel(
            this.props.id,
            data.cmd,
            data.value,
            data.r,
            data.g,
            data.b,
            function(err, data) {
                callback(err, data);
            });
    },

    render: function() {
        var sliderCmp;
        if (this.props.output === 'continuous') {
            sliderCmp = (
                <div>
                    <div className="slider"></div>
                    <span className={ClassNames({
                        "value": true,
                        "hidden": this.state.value === -1})}>{this.state.value}%</span>
                </div>
            );
        }

        var rgbCmp;
        if (this.props.output === 'rgb') {
            rgbCmp = (
                <div className="rgbWrapper"></div>
            );
        }

        var classes = { 'cmp-ZoneControl': true };
        classes[this.props.output] = true;
        return (
            <div className={ClassNames(classes)}>
                <div className="clearfix">
                    <div className="name pull-left">
                        {this.props.name}
                    </div>
                    <div className="onOffWrapper pull-right">
                        <input
                            className="switch-indeterminate"
                            type="checkbox"
                            defaultChecked={true}
                            data-indeterminate="true"></input>
                    </div>
                </div>
                {sliderCmp}
                {rgbCmp}
            </div>
        );
    }
});
module.exports = ZoneControl;
