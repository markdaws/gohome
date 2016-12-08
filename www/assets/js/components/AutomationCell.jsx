var React = require('react');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'AutomationCell',
    prefix: 'b-'
});
require('../../css/components/AutomationCell.less')

var AutomationCell = React.createClass({
    render: function() {
        return (
            <div {...classes()}>
                <div {...classes('icon')}>
                    <i className="icon ion-ios-cog-outline"></i>
                </div>
                <div {...classes('name')}>
                    {this.props.automation.name}
                </div>
            </div>
        );
    }
});
module.exports = AutomationCell;
