var React = require('react');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'Spinner',
    prefix: 'b-'
});
require('../../css/components/Spinner.less')

var Spinner = React.createClass({
    render: function() {
        return (
            <div {...classes('', this.props.hidden ? 'hidden' : '')}>
                <div {...classes('icon')}></div>
            </div>
        );
    }
});
module.exports = Spinner;
