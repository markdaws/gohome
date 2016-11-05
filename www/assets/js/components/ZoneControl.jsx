var ClassNames = require('classnames');
var React = require('react');
var ReactDOM = require('react-dom');
var CssMixin = require('./CssMixin.jsx');
var Api = require('../utils/API.js');
var ClassNames = require('classnames');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'ZoneControl',
    prefix: 'b-'
});
require('../../css/components/ZoneControl.less')

var ZoneControl = React.createClass({
    mixins: [CssMixin],
    getInitialState: function() {
        var level = this.props.level;
        if (!level) {
            level = {
                value: -1,
                r: 0,
                g: 0,
                b: 0
            };
        }
        return {
            level: level
        }
    },

    initSlider: function() {
        var sliders = $(ReactDOM.findDOMNode(this)).find('.b-ZoneControl__slider');
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
                step: 1,
                orientation: 'horizontal',
                range: {
                    min: 0,
                    max: 100
                }
            });
        slider.noUiSlider.set(0);
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

        return sw;
    },

    initRGB: function() {
        var wrapper = $(ReactDOM.findDOMNode(this)).find('.rgbWrapper');
        if (!wrapper || wrapper.length === 0) {
            return;
        }

        var c = new Colors({
            'color': 'rgba(255,0,0,1.0)'
            });
        wrapper.colorpicker({
            format: 'rgb',
            container: true,
            inline: true,
            color: c,
            buildCallback: function($elm) {
                debugger;
            },
            renderCallback: function() {
                debugger;
            }
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

        return wrapper;
    },

    makeLevel: function(val, r, g, b) {
        return {
            value: val,
            r: r,
            g: g,
            b: b
        }
    },
    
    sliderChanged: function(slider) {
        this.setState({ level: this.makeLevel(parseInt(slider.get(), 10), 0, 0, 0) });
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
        this._slider = this.initSlider();
        this._switch = this.initSwitch(this._slider);
        this._colorPicker = this.initRGB();

        this.props.didMount && this.props.didMount(this);
    },

    componentWillUnmount: function() {
        this.props.willUnmount && this.props.willUnmount();
    },

    setValue: function(cmd, value, r, g, b, callback) {
        this.setState({ level: this.makeLevel(value, r, g, b) });
        this.send({
            cmd: cmd,
            value: parseFloat(value),
            r: r,
            g: g,
            b: b
        }, callback);
    },

    toggleOn: function(slider) {
        var cmd, targetValue;
        
        if (this.state.level.value !== 0 || this.state.level.r !== 0 || this.state.level.g !== 0 || this.state.level.b !== 0) {
            cmd = 'turnOff';
            targetValue = 0;
        } else {
            cmd = 'turnOn';
            targetValue = 100;
        }
        
        slider && slider.set(targetValue);
        this.setValue(cmd, targetValue, 0, 0, 0, function(err) {
            if (err) {
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

    monitorData: function(data) {
        if (!data || !data.zones) {
            return;
        }
        var val = data.zones[this.props.id];
        if (val == undefined) {
            return;
        }

        this.setState({ level: this.makeLevel(Math.round(val.value), 0, 0, 0) });

        switch(this.props.output) {
            case 'continuous':
                this._slider && this._slider.set(Math.round(val.value));
                break;

            case 'rgb':
                //TODO: color picker not changing
                /*
                this._colorPicker.colorpicker({
                    color: new Colors({ color: 'rgba(' + val.r + ',' + val.g + ',' + val.b + ',1.0)' })
                });*/
                break;
        }
        this._switch && this._switch.bootstrapSwitch('state', val.value > 0, true);
    },

    render: function() {
        var sliderCmp;
        if (this.props.output === 'continuous') {
            sliderCmp = (
                <div>
                    <div {...classes('slider')}></div>
                    <span {...classes('value', this.state.level.value === -1 ? 'hidden' : '')}>{this.state.level.value}%</span>
                </div>
            );
        }

        var rgbCmp;
        if (this.props.output === 'rgb') {
            rgbCmp = (
                <div {...classes('rgb')}></div>
            );
        }

        return (
            <div {...classes('', this.props.output)}>
                <div className="clearfix">
                    <div {...classes('name', '', 'pull-left')}>
                        {this.props.name}
                    </div>
                    <div {...classes('on-off', '', 'pull-right')}>
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
