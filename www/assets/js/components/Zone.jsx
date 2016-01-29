var React = require('react');
var ReactDOM = require('react-dom');
var CssMixin = require('./CssMixin.jsx');

module.exports = React.createClass({
    mixins: [CssMixin],
    getInitialState: function() {
        return {
            value: 0,
            showSlider: false,
            slider: null
        }
    },

    componentDidMount: function() {
        var self = this;
        
        switch (this.props.output) {
        case 'binary':
        case 'continuous':
            var s = $(ReactDOM.findDOMNode(this)).find('.valueSlider');
            s.slider({ reversed: false });
            self.setState({ slider: s });
            s.on('change', function(evt) {
                self.setState({ value: evt.value.newValue });
            });
            s.on('slideStop', function(evt) {
                self.setValue('setLevel', evt.value, 0, 0, 0, function(err) {
                    if (err) {
                        //TODO:
                        console.error(err);
                    }
                });
                return false;
            });
            break;

        case 'rgb':
            var $el = $(ReactDOM.findDOMNode(this)).find('.zone-rgb .clickInfo span')
            $el.colorPicker({
                doRender:false,
                opacity: false,
                margin: '0px 0px 0px -30px',
                renderCallback: function($e, toggled) {
                    if (toggled !== undefined) {
                        // only send a value when the user actually interacts with the
                        // control not when it is first shown/hidden
                        return;
                    }
                    var rgb = this.color.colors.rgb;
                    self.setValue(
                        'setLevel',
                        0,
                        parseInt(rgb.r * 255),
                        parseInt(rgb.g * 255),
                        parseInt(rgb.b * 255),
                        function(err) {
                            if (err) {
                                console.error(err);
                            }
                        }
                    );
                }
            });
            break;
        }
    },

    infoClicked: function(evt) {
        evt.stopPropagation();
        evt.preventDefault();

        if (!this.isRgb()) {
            this.setState({ showSlider: true });
        }
    },

    isRgb: function() {
        return this.props.output === 'rgb';
    },
    
    setValue: function(cmd, value, r, g, b, callback) {
        if (!this.isRgb()) {
            this.state.slider.slider('setValue', value, false, true);
        }
        //TODO: Need rgb
        this.setState({ value: value });
        this.send({
            cmd: cmd,
            value: parseFloat(value),
            r: r,
            g: g,
            b: b
        }, callback);
    },

    toggleOn: function(evt) {
        evt.stopPropagation();
        evt.preventDefault();

        var cmd, level;
        if (this.state.value !== 0) {
            cmd = 'turnOff';
            level = 0;
        } else {
            cmd = 'turnOn';
            level = 100;
        }
        this.setValue(cmd, level, 0, 0, 0, function(err) {
            if (err) {
                console.error(err);
            }
        });
    },

    send: function(data, callback) {
        $.ajax({
            url: '/api/v1/systems/1/zones/' + this.props.id,
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(data),
            success: function(data) {
                callback();
            },
            error: function(xhr, status, err) {
                callback(err);
            }
        });
    },

    render: function() {
        var value = this.state.value;
        var icon = this.props.type === 'light' ? 'fa fa-lightbulb-o' : 'fa fa-picture-o';

        var stepSize
        switch (this.props.output) {
        case 'continuous':
            stepSize = 1;
            break;
        case 'binary':
            stepSize = 100;
            break;
        case 'rgb':
            break;
        default:
            stepSize = 1;
        }
        
        return (
            <div className="cmp-Zone col-xs-12 col-sm-4 col-md-4 col-lg-4 clearfix">
              <button className={"btn btn-primary zone" + (this.isRgb() ? " zone-rgb" : "")}>
                <i className={icon}></i>
                <span className="name">{this.props.name}</span>
                <div className={"sliderWrapper pull-right" + (this.state.showSlider ? "" : " hidden")} >
                  <input className="valueSlider" type="text" data-slider-value="0" data-slider-min="00" data-slider-max="100" data-slider-step={stepSize} data-slider-orientation="horizontal"></input>
                  <span className="level pull-right">{this.state.value}%</span>
                </div>
                <div className="clearfix footer">
                  <div className={"clickInfo pull-right" + (this.state.showSlider ? " hidden" : "")}>
                    <span onClick={this.infoClicked}>Click to control</span>
                  </div>
                  <a className="btn btn-link btnToggle pull-left" onClick={this.toggleOn}>
                    <i className="fa fa-power-off"></i>
                  </a>
                </div>
              </button>
            </div>
        )
    }
});
