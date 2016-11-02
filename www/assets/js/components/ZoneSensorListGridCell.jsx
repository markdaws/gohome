var React = require('react');
var ReactDOM = require('react-dom');
var Classnames = require('classnames');

var ZoneSensorListGridCell = React.createClass({
    getInitialState: function() {
        return {
            level: this.props.level || { value: 0, r:0, g:0, b:0 },
            attr: null
        };
    },

    setLevel: function(level) {
        this.setState({ level: level });

        /*
        var $this = $(ReactDOM.findDOMNode(this));
        //TODO: Clip rect is not updating on ios safari, so just change the
        //opacity of the bulb for now
        //var y = parseInt(30 + (55 * ((100-val)/100)));
        //$this.find('.clipRect').attr('y', y);
        */
        //$this.find('.light').css('opacity', val/100);
    },

    setAttr: function(attr) {
        this.setState({ attr: attr });
    },

    shouldComponentUpdate: function(nextProps, nextState) {
        //TODO: Fix
        return true;
        if (nextProps.zone && this.props.zone && (this.props.zone.name !== nextProps.zone.name)) {
            return true;
        }
        if (nextProps.sensor && this.props.sensor && (this.props.sensor.name !== nextProps.sensor.name)) {
            return true;
        }
        if (nextState.level && this.state.level && (nextState.level.value !== this.state.level.value)) {
            //TODO: RGB
            return true;
        }
        if (nextState.attr && this.state.attr && (nextState.attr.value !== this.state.attr.value)) {
            //TODO: RGB
            return true;
        }

        return false;
    },
    
    render: function() {
        var icon1, icon2, name;
        var type;
        var val = '';
        var opacity = 0;
        var color = 'yellow';
        
        if (this.props.zone) {
            switch (this.props.zone.type) {
                case 'light':
                    type = 'light';
                    icon1 = 'icon ion-ios-lightbulb-outline';
                    break;
                case 'shade':
                    type = 'shade';
                    icon1 = 'icon ion-ios-arrow-thin-up';
                    icon2 = 'icon ion-ios-arrow-thin-down';
                    break;
                case 'switch':
                    type = 'switch';
                    icon1 = 'icon ion-ios-bolt-outline';
                    break;
                default:
                    icon1 = 'icon ion-ios-help-empty';
                    break;
            }
            name = this.props.zone.name;

            if (this.state.level) {
                if (this.props.zone.output === 'binary') {
                    opacity = this.state.level.value === 0 ? 0 : 1;
                    val = this.state.level.value === 0 ? 'off' : 'on';
                } else if (this.props.zone.output === 'rgb') {
                    opacity = 1;

                    var lev = this.state.level;
                    val = this.state.level.value === 0 ? 'off' : 'on';
                    color = "#" + ((1 << 24) + (lev.r << 16) + (lev.g << 8) + lev.b).toString(16).slice(1);
                } else {
                    opacity = this.state.level.value / 100;
                    val = this.state.level.value + '%';
                }
                
                if (this.props.zone && this.props.zone.type === 'switch') {
                    if (this.state.level.value === 0) {
                        val = 'off';
                    } else {
                        val = 'on';
                    }
                }
            }
        } else {
            icon1 = 'icon ion-ios-pulse';
            type = 'sensor';
            name = this.props.sensor.name;

            if (this.state.attr) {
                val = this.state.attr.value;

                // If there is a states map, which gives value -> ui string then
                // use that string instead of the raw value
                var uiVal = this.state.attr.states && this.state.attr.states[val];
                if (uiVal) {
                    val = uiVal;
                }
            }
        }

        var icon1Cmp, icon2Cmp;
        icon1Cmp = <i className={icon1}></i>;
        if (icon2) {
            icon2Cmp = <i className={icon2}></i>;
        }

        var typeClass = {};
        typeClass[type] = true;
        return (
        <div className={Classnames("cmp-ZoneSensorListGridCell", typeClass)}>
            <div className="icon">
                {icon1Cmp}
                {icon2Cmp}
            </div>
            <svg
                viewBox="0 0 200 200"
                xmlns="http://www.w3.org/2000/svg"
                xlinkHref="http://www.w3.org/1999/xlink">
                <g>
                    <clipPath id="lightClip">
                        <rect className="clipRect" x="0" y="30" width="200" height="65" />
                    </clipPath>
                </g>
                <path
                    className="switch"
                    d="M105 45 L82 75 L100 75 L95 100 L120 67 L100 65"
                    stroke="yellow"
                    style={{'opacity': opacity}}
                    fill="yellow"></path>
                <circle
                    className="light"
                    cx="100"
                    cy="55"
                    r="25"
                    fill={color}
                    clipPath="url(#lightClip)"
                    style={{'opacity': opacity, 'clipPath':'url(#lightClip)'}}/>
            </svg>
            <div className="level">
                {val}
            </div>
            <div className="name">
                {name}
            </div>
        </div>
        );
    }
});
module.exports = ZoneSensorListGridCell;
