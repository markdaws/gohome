var React = require('react');

var Scene = React.createClass({
    handleClick: function(event) {
        $.ajax({
            url: '/api/v1/systems/1/scenes/active',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify({ id: this.props.scene.id }),
            success: function(data) {
                //TODO: Common way in UI to display success/error
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err.toString());
            }.bind(this)
        });
    },

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
