var React = require('react');
var ReactDOM = require('react-dom');
var Classnames = require('classnames');

var ZoneSensorListGridCell = React.createClass({
    getInitialState: function() {
        return {
            level: this.props.level || { value: 0, r:0, g:0, b:0 }
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

    shouldComponentUpdate: function(nextProps, nextState) {
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
        return false;
    },
    
    render: function() {
        var icon1, icon2, name;
        var type;
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
        } else {
            icon1 = 'icon ion-ios-pulse';
            type = 'sensor';
            name = this.props.sensor.name;
        }

        var icon1Cmp, icon2Cmp;
        icon1Cmp = <i className={icon1}></i>
        if (icon2) {
            icon2Cmp = <i className={icon2}></i>
        }

        var val = '';
        var opacity = 0;
        if (this.state.level) {
            opacity = this.state.level.value / 100;
            val = this.state.level.value + '%';
        
            if (this.props.zone && this.props.zone.type === 'switch') {
                if (this.state.level.value === 0) {
                    val = 'off';
                } else {
                    val = 'on';
                }
            }
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
                <path className="switch" d="M105 45 L82 75 L100 75 L95 100 L120 67 L100 65" stroke="yellow" fill="yellow"></path>
                <circle
                    className="light"
                    cx="100"
                    cy="55"
                    r="25"
                    fill="yellow"
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
