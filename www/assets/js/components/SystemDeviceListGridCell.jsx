var React = require('react');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SystemDeviceListGridCell',
    prefix: 'b-'
});
require('../../css/components/SystemDeviceListGridCell.less')

var SystemDeviceListGridCell = React.createClass({
    render: function() {
        return (
            <div {...classes()}>
                <div {...classes('icon')}>
                    <i className="icon ion-cube"></i>
                </div>
                <div {...classes('name')}>
                    {this.props.device.name}
                </div>
            </div>
        );
    }
});
module.exports = SystemDeviceListGridCell;
