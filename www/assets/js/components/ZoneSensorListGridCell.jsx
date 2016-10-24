var React = require('react');

var ZoneSensorListGridCell = React.createClass({
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
        return (
            <div className="cmp-ZoneSensorListGridCell">
                <div className="icon">
                    {icon1Cmp}
                    {icon2Cmp}
                </div>
                <div className="name">
                    {name}
                </div>
            </div>
        );
    }
});
module.exports = ZoneSensorListGridCell;
