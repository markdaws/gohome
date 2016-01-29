var React = require('react');
var Action = require('./Action.jsx');

module.exports = React.createClass({
    handleClick: function(action) {
        this.props.selected(action);
    },

    render: function() {
        var self = this;
        var actionNodes = this.props.actions.map(function(action) {
            return (
                <Action data={action} selected={self.handleClick} key={action.name}/>
            );
        });
        return (
            <div className="cmp-ActionList clearfix">
              {actionNodes}
            </div>
        );
    }
});
