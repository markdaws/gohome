var React = require('react');
var Api = require('../utils/API.js')

var Scene = React.createClass({
    handleClick: function(event) {
        Api.sceneActivate(this.props.scene.id, function(err, data) {
            //TODO: Show error/success
        });
z    },

    render: function() {
        return (
            <div className="cmp-Scene col-xs-6 col-sm-3 col-md-3 col-lg-3">
              <a role="button" className="btn btn-primary scene" onClick={this.handleClick}>
                <div>
                  <i className="fa fa-sliders"></i>
                </div>
                <span className="name">{this.props.scene.name}</span>
              </a>
            </div>
        )
    }
});
module.exports = Scene;
