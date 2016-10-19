var React = require('react');
var Api = require('../utils/API.js')

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
            <div className="cmp-SceneControl">
                <div className="name">
                    {this.props.scene.name}
                </div>
                <div className="activateWrapper">
                    <a role="button" className="btn btn-primary scene" onClick={this.handleClick}>
                        <span className="name">Activate</span>
                    </a>
                </div>
            </div>
        )
    }
});
module.exports = SceneControl;
