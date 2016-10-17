var React = require('react');

var ZoneListGridCell = React.createClass({
    render: function() {
        var icon1, icon2;
        switch (this.props.zone.type) {
            case 'light':
                icon1 = 'icon ion-ios-lightbulb-outline';
                break;
            case 'shade':
                icon1 = 'icon ion-ios-arrow-thin-up';
                icon2 = 'icon ion-ios-arrow-thin-down';
                break;
            case 'outlet':
                icon1 = 'icon ion-ios-bolt-outline';
                break;
            default:
                icon1 = 'icon ion-ios-help-empty';
                break;
        }

        var icon1Cmp, icon2Cmp;
        icon1Cmp = <i className={icon1}></i>
        if (icon2) {
            icon2Cmp = <i className={icon2}></i>
        }
        return (
            <div className="cmp-ZoneListGridCell">
                <div className="icon">
                    {icon1Cmp}
                    {icon2Cmp}
                </div>
                <div className="name">
                    {this.props.zone.name}
                </div>
            </div>
        );
    }
});
module.exports = ZoneListGridCell;
