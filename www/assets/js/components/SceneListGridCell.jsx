var React = require('react');

var SceneListGridCell = React.createClass({
    render: function() {
        return (
            <div className="cmp-SceneListGridCell">
                <div className="icon">
                    <i className="icon ion-ios-settings"></i>
                </div>
                <div className="name">
                    {this.props.scene.name}
                </div>
            </div>
        );
    }
});
module.exports = SceneListGridCell;
