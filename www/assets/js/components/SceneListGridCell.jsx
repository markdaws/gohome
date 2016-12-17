var React = require('react');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SceneListGridCell',
    prefix: 'b-'
});
require('../../css/components/SceneListGridCell.less')

var SceneListGridCell = React.createClass({
    render: function() {
        return (
            <div {...classes()}>
                <div {...classes('icon')}>
                    <i className="icomoon-ion-ios-settings"></i>
                </div>
                <div {...classes('name')}>
                    {this.props.scene.name}
                </div>
            </div>
        );
    }
});
module.exports = SceneListGridCell;
