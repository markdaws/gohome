var React = require('react');

var Trigger = React.createClass({
    handleClick: function(evt) {
        evt.preventDefault();
        evt.stopPropagation();
        this.props.selected(this.props.data);
    },

    render: function() {
        return (
            <div className="cmp-Trigger pull-left">
              <button className="btn btn-primary" onClick={this.handleClick}>
                <h4>{this.props.data.name}</h4>
                <p>{this.props.data.description}</p>
              </button>
            </div>
        );
    }
});
module.exports = Trigger;