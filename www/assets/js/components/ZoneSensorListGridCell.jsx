var React = require('react');
var ReactDOM = require('react-dom');

var ZoneSensorListGridCell = React.createClass({
    getInitialState: function() {
        return {
            level: this.props.level
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
        if (nextState.level && (nextState.level !== this.state.level)) {
            return true;
        }
        return false;
    },
    
    render: function() {
        var icon1, icon2, name;

        if (this.props.zone) {
            switch (this.props.zone.type) {
                case 'light':
                    icon1 = 'icon ion-ios-lightbulb-outline';
                    break;
                case 'shade':
                    icon1 = 'icon ion-ios-arrow-thin-up';
                    icon2 = 'icon ion-ios-arrow-thin-down';
                    break;
                case 'switch':
                    icon1 = 'icon ion-ios-bolt-outline';
                    break;
                default:
                    icon1 = 'icon ion-ios-help-empty';
                    break;
            }
            name = this.props.zone.name;
        } else {
            icon1 = 'icon ion-ios-pulse';
            name = this.props.sensor.name;
        }

        var icon1Cmp, icon2Cmp;
        icon1Cmp = <i className={icon1}></i>
        if (icon2) {
            icon2Cmp = <i className={icon2}></i>
        }

        var val = 0;
        if (this.state.level) {
            val = this.state.level.value / 100;
        }
        return (
        <div className="cmp-ZoneSensorListGridCell">
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
                <circle
                    className="light"
                    cx="100"
                    cy="55"
                    r="25"
                    fill="yellow"
                    clipPath="url(#lightClip)"
                    style={{'opacity': val, 'clipPath':'url(#lightClip)'}}/>
            </svg>
            <div className="name">
                {val}{name}
            </div>
        </div>
        );
    }
});
module.exports = ZoneSensorListGridCell;
