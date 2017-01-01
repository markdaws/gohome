var React = require('react');
var Api = require('../utils/API.js')
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SceneControl',
    prefix: 'b-'
});
require('../../css/components/SceneControl.less')

var SceneControl = React.createClass({
    handleClick: function(event) {
        Api.sceneActivate(this.props.scene.id, function(err, data) {
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
                    {this.props.scene.name}
                </div>
                <a role="button" {...classes('activate', '', 'btn btn-primary')} onClick={this.handleClick}>
                    Activate
                </a>
            </div>
        )
    }
});
module.exports = SceneControl;
