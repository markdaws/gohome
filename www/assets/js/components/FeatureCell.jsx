var React = require('react');
var Feature = require('../feature.js');
var Attribute = require('../attribute.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'FeatureCell',
    prefix: 'b-'
});
require('../../css/components/FeatureCell.less')

var FeatureCell = React.createClass({
    getDefaultProps: function() {
        return {
            hasError: false,
            hasSuccess: false,
            showCheckbox: false,
            //TODO: Rename
            showLevel: true,
            checkboxChecked: true
        }
    },

    getInitialState: function() {
        return {
            attrs: this.props.feature.attrs,
            checkboxChecked: this.props.checkboxChecked
        };
    },

    isChecked: function() {
        return this.state.checkboxChecked;
    },

    checkboxClicked: function(evt) {
        evt.stopPropagation();
    },

    checkboxChanged: function(evt) {
        this.setState({checkboxChecked: evt.target.checked});
        this.props.chkBxChanged && this.props.chkBxChanged(this.props.id, evt.target.checked);
    },

    setAttrs: function(attrs) {
        this.setState({ attrs: attrs });
    },

    shouldComponentUpdate: function(nextProps, nextState) {
        //TODO: Fix
        return true;
    },

    render: function() {
        var icon1, icon2;

        var name = this.props.feature.name;
        var attrs = this.state.attrs;
        var val = '';
        var opacity = 0;
        var color = 'yellow';
        switch(this.props.feature.type) {
            case Feature.Type.LightZone:
                icon1 = 'icon ion-ios-lightbulb-outline';

                var isOff = true
                var onOffVal = attrs[Feature.LightZone.AttrIDs.OnOff].value;
                if (!onOffVal) {
                    val = '';
                } else {
                    if (onOffVal === 2) {
                        isOff = false
                        opacity = 1;
                    }
                    val = Attribute.OnOff.States[onOffVal];
                }

                if (!isOff && attrs[Feature.LightZone.AttrIDs.HSL]) {
                    color = attrs[Feature.LightZone.AttrIDs.HSL].value;
                    if (color == null) {
                        opacity = 0;
                    } else {
                        opacity = 1;
                    }
                }

                if (!isOff && attrs[Feature.LightZone.AttrIDs.Brightness]) {
                    // The light zone supports brightness, show the current intensity
                    val = attrs[Feature.LightZone.AttrIDs.Brightness].value;
                    if (val == null) {
                        val = '';
                    } else {
                        opacity = val / 100;
                        val = parseInt(val, 10) + '%';
                    }
                }

                break;

            case Feature.Type.Switch:
                icon1 = 'icon ion-ios-bolt-outline';
                var onOffVal = attrs[Feature.Switch.AttrIDs.OnOff].value;
                if (!onOffVal) {
                    val = '';
                } else {
                    if (onOffVal === 2) {
                        opacity = 1;
                    }
                    val = Attribute.OnOff.States[onOffVal];
                }
                break;

            case Feature.Type.Outlet:
                icon1 = 'icon ion-outlet';
                var onOffVal = attrs[Feature.Outlet.AttrIDs.OnOff].value;
                if (!onOffVal) {
                    val = '';
                } else {
                    if (onOffVal === 2) {
                        opacity = 1;
                    }
                    val = Attribute.OnOff.States[onOffVal];
                }
                break;

            case Feature.Type.WindowTreatment:
                icon1 = 'icon ion-ios-arrow-thin-up';
                icon2 = 'icon ion-ios-arrow-thin-down';
                val = attrs[Feature.WindowTreatment.AttrIDs.Offset].value;
                if (val == null) {
                    val = '';
                } else if (val === 0) {
                    val = 'Closed';
                } else if (val === 100) {
                    val = 'Open';
                } else {
                    val = val + '%';
                }
                break;

            case Feature.Type.HeatZone:
                icon1 = 'icon ion-ios-flame-outline';
                var current = attrs[Feature.HeatZone.AttrIDs.CurrentTemp].value;
                var target = attrs[Feature.HeatZone.AttrIDs.TargetTemp].value;

                if (current == null || target == null) {
                    val = '';
                } else if (current === target) {
                    val = current + '°F';
                } else {
                    val = current + '°F → ' + target + '°F';
                }
                break;

            case Feature.Type.Sensor:
                icon1 = 'icon ion-ios-pulse';

                // Each sensor has only one attribute, pick it out
                var attribute = attrs[Object.keys(attrs)[0]];
                if (!attribute.value) {
                    val = '';
                } else {
                    val = Attribute[attribute.type].States[attribute.value];
                }
                break;

            default:
                icon1 = 'icon ion-ios-help-empty';
        }

        var icon1Cmp, icon2Cmp;
        icon1Cmp = <i className={icon1}></i>;
        if (icon2) {
            icon2Cmp = <i className={icon2}></i>;
        }

        if (!this.props.showLevel) {
            val = null;
        }

        var chkBx;
        if (this.props.showCheckbox) {
            chkBx = (
                <input
                    {...classes('checkbox')}
                    type="checkbox"
                    onChange={this.checkboxChanged}
                    onClick={this.checkboxClicked}
                    checked={this.state.checkboxChecked}
                ></input>
            );
        }

        var state = '';
        if (this.props.hasError) {
            state = 'error';
        } else if (this.props.hasSuccess) {
            state = 'success';
        }
        return (
            <div {...classes('', state)}>
                {chkBx}
                <div {...classes('icon')}>
                    {icon1Cmp}
                    {icon2Cmp}
                </div>
                <svg
                    {...classes('switch-color', this.props.feature.type)}
                    viewBox="0 0 200 200"
                    xmlns="http://www.w3.org/2000/svg"
                    xlinkHref="http://www.w3.org/1999/xlink">
                    <path
                        d="M105 45 L82 75 L100 75 L95 100 L120 67 L100 65"
                        stroke="yellow"
                        style={{'opacity': opacity}}
                        fill="yellow"
                        transform="translate(-1, -10)"
                    ></path>
                </svg>
                <svg
                    {...classes('light-color', this.props.feature.type)}
                    viewBox="0 0 200 200"
                    xmlns="http://www.w3.org/2000/svg"
                    xlinkHref="http://www.w3.org/1999/xlink">
                    <g>
                        <clipPath id="lightClip">
                            <rect className="clipRect" x="0" y="20" width="200" height="65" />
                        </clipPath>
                    </g>
                    <circle
                        cx="100"
                        cy="48"
                        r="22"
                        fill={color}
                        clipPath="url(#lightClip)"
                        style={{'opacity': opacity, 'clipPath':'url(#lightClip)'}}/>
                </svg>
                <div {...classes('level')}>
                    {val}
                </div>
                <div {...classes('name')}>
                    {name}
                </div>
            </div>
        );
    }
});
module.exports = FeatureCell;
