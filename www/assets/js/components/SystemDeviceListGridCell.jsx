var React = require('react');

var SystemDeviceListGridCell = React.createClass({
    render: function() {
        return (
            <div className="cmp-SystemDeviceListGridCell">
                <div className="icon">
                    <i className="icon ion-cube"></i>
                </div>
                <div className="name">
                    {this.props.device.name}
                </div>
            </div>
        );
    }
});
module.exports = SystemDeviceListGridCell;
