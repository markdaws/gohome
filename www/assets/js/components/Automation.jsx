var React = require('react');
var Api = require('../utils/API.js')
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'Automation',
    prefix: 'b-'
});
require('../../css/components/Automation.less')

var Automation = React.createClass({
    handleClick: function(event) {
        Api.automationTest(this.props.automation.id, function(err, data) {
            if (err) {
                //TODO: Show error/success
                console.error(err);
            }
        });
    },

    render: function() {
        return (
            <div {...classes()}>
                <div {...classes('name')}>
                    {this.props.automation.name}
                </div>
                <a role="button" {...classes('activate', '', 'btn btn-primary')} onClick={this.handleClick}>
                    Test
                </a>
            </div>
        )
    }
});
module.exports = Automation;
